package TaskSyncService

import (
	"bytes"
	"crypto/tls"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"net/mail"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bytbox/go-pop3"
	"nexsoft.co.id/tes-worker/config"
	"nexsoft.co.id/tes-worker/constant"
	"nexsoft.co.id/tes-worker/dao"
	"nexsoft.co.id/tes-worker/model/backgroundJobModel"
	"nexsoft.co.id/tes-worker/model/errorModel"
	"nexsoft.co.id/tes-worker/repository"
	"nexsoft.co.id/tes-worker/serverConfig"

	gomail "github.com/veqryn/go-email/email"
)

var (
	client         *pop3.Client
	listExtenstion string
)

// GetSyncElasticChildTask
func (input taskService) GetSyncChildGetEmailTask() backgroundJobModel.ChildTask {
	return backgroundJobModel.ChildTask{
		Group: constant.JobProcessSynchronizeGroup,
		Type:  constant.JobProcessGetEmailType,
		Name:  constant.JobProcessSyncTaskGetEmail,
		Data: backgroundJobModel.BackgroundServiceModel{
			SearchByParam: nil,
			IsCheckStatus: false,
			CreatedBy:     0,
			Data:          nil,
		},
		GetCountData: dao.EmailDAO.GetCountEmail,
		DoJob:        input.GetAllEmail,
	}
}

func (input taskService) GetAllEmail(db *sql.DB, _ interface{}, childJob *repository.JobProcessModel) (err errorModel.ErrorModel) {
	fmt.Println("\n ==================")
	fmt.Println("  TASK Get Email   !!  ")
	listExtenstion = ".xls, .xlsx, .doc, .docx, .pdf, .jpeg, .zip, .rar"
	currentTime := time.Now()

	fmt.Println("Start : ", currentTime.String())
	readAllEmail()
	fmt.Println("End   : ", currentTime.String())
	// Set status Done
	fmt.Println("\n ==================")
	return
}

func readAllEmail() {
	fmt.Println(" ----> READ EMAIL ")
	// client := serverConfig.ServerAttribute.Client

	username := config.ApplicationConfiguration.GetEmailUsername()
	password := config.ApplicationConfiguration.GetEmailPassword()
	serverAddress := config.ApplicationConfiguration.GetEmailAddress()
	secure := config.ApplicationConfiguration.GetEmailSecure()

	fmt.Println("Connect to ", serverAddress, " user ", username, " pass ", password, " secure ", secure)

	var dialErr error
	if secure {
		conn, err := tls.Dial("tcp", serverAddress, &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			fmt.Println("error Dial ", err.Error())
			return
		}
		client, dialErr = pop3.NewClient(conn)

	} else {
		client, dialErr = pop3.Dial(serverAddress)
	}

	if dialErr != nil {
		fmt.Println("error Dial ", dialErr.Error())
		return
	}

	authErr := client.Auth(username, password)
	if authErr != nil {
		fmt.Println("error Auth dial ", authErr.Error())
		return
	}

	mailCount, mailBoxSize, err := client.Stat()
	fmt.Println("MailCount: ", mailCount, " MailBoxSize:", uint64(mailBoxSize))
	if err != nil {
		fmt.Println("Err Retrieve  ")
		return
	}

	msgs, _, err := client.ListAll()
	if err != nil {
		fmt.Println("Err get list all ")
		return
	}

	fmt.Println("\nmessage START =================================================== ")
	for i := 0; i < len(msgs); i++ {
		downloadEMailByMessageID(client, msgs[i])
		if err != nil {
			fmt.Println("Download mail failure, Err: ", err)
		}
	}
	fmt.Println("\nmessage FINISH =================================================== ")
}

func downloadEMailByMessageID(client *pop3.Client, index int) {

	fmt.Println("Downloading idx =", index)
	var mailContent string
	var err error
	mailContent, err = client.Retr(index)
	if err != nil {
		fmt.Println("err Retr=> ", err)
		return
	}
	var lastNewline = 0
	for i, b := range mailContent {
		if b != '\n' {
			lastNewline = i
			break
		}
	}
	rawMailContent := []byte(mailContent)[lastNewline:]
	var emailModel repository.TicketEmail
	var emailAttach repository.TicketEmailAttachment

	readEmail(rawMailContent, &emailModel)
	parsingEmail(rawMailContent, index, &emailModel, &emailAttach)

	saveToDatabase(emailModel, emailAttach)
}

func readEmail(rawMailContent []byte, emailModel *repository.TicketEmail) error {

	fmt.Println("READ EMAIL...")
	var message *mail.Message
	message, err := mail.ReadMessage(bytes.NewReader(rawMailContent))
	if err != nil {
		return err
	}

	address, err := mail.ParseAddress(message.Header.Get("From"))
	if err != nil {
		fmt.Println("parse from address failure, Err:", err)
		return err
	}
	// fmt.Println("subject=>", message.Header.Get("Subject"))

	// fmt.Println("Date=>", message.Header.Get("Date"))
	// fmt.Println("Content-Type=>", message.Header.Get("Content-Type"))
	// fmt.Println("Delivered-To=>", message.Header.Get("Delivered-To"))
	// fmt.Println("Message-ID=>", message.Header.Get("Message-ID"))
	// fmt.Println("address name=>", address.Name)
	// fmt.Println("address email =>", address.Address)

	layout := "2006-01-02T15:04:05.000Z"
	emailDate, _ := time.Parse(layout, message.Header.Get("Date"))

	emailModel.TicketID.Int64 = 0
	emailModel.MessageID.String = message.Header.Get("Message-ID")
	emailModel.EmailSubject.String = message.Header.Get("Subject")
	emailModel.EmailDate.Time = emailDate
	emailModel.EmailSender.String = address.Address
	emailModel.EmailSenderName.String = address.Name
	// emailModel.TextPlain
	// emailModel.TextHTML

	fmt.Println("\n \n  ")
	return nil
}

func parsingEmail(rawMailContent []byte, index int, emailModel *repository.TicketEmail, attach *repository.TicketEmailAttachment) {

	var fileName string
	var pathsave string
	msg, err := gomail.ParseMessage(bytes.NewReader(rawMailContent))

	if err != nil {
		fmt.Println("Err parse from gomail ", err.Error())
		return
	}
	for _, part := range msg.MessagesAll() {

		fileName = ""
		mediaType, _, err := part.Header.ContentType()
		_, paramAttach, _ := part.Header.ContentDisposition()
		// fmt.Println("Param attach ==>", paramAttach)
		fileName = paramAttach["filename"]

		if err != nil {
			fmt.Println("Err from gomail ", err)
			break
		}
		switch mediaType {
		case "image/jpeg":
			fmt.Println("image/jpeg = FOUND = ", fileName)
		case "image/jpg":
			fmt.Println("image/jpg = FOUND = ", fileName)
		case "text/plain":
			fmt.Println("textplain = FOUND = ", fileName)
			if fileName == "" {
				// fmt.Println("Body textplain ==>", string(part.Body))
				emailModel.TextPlain.String = string(part.Body)
			}
		case "text/html":
			fmt.Println("text html = FOUND = ", fileName)
			if fileName == "" {
				// fmt.Println("Body HTML ==>", string(part.Body))
			}
		case "application/pdf":
			fmt.Println("pdf = FOUND = ", fileName)
		default:
			fmt.Println("default => ", mediaType, " = FOUND = ", fileName)
		}

		if fileName != "" {
			validExt := validateAllowExtention(fileName)
			if validExt {
				// save attachment
				attach.Filename.String = fileName
				attach.Filetype.String = mediaType

				pathsave = saveAttachment(fileName, part.Body, index)
			}

		}
	}

	if pathsave == "" {
		pathsave = generatePathAttachment(index)
	}
	// save RAW email
	filePathAttachment := filepath.Join(pathsave, strconv.Itoa(index))
	ioutil.WriteFile(filePathAttachment, rawMailContent, 0644)

}

func validateAllowExtention(fileName string) bool {

	fmt.Println("Filename : ", fileName)
	fileExtension := path.Ext(fileName)
	fmt.Println("Is enable extenstion ==> [", fileExtension, "] ")

	allowExtenstions := strings.Split(listExtenstion, ",")
	fmt.Println("Allowed extenstion ==> ", allowExtenstions)

	for _, allowExtenstion := range allowExtenstions {
		fmt.Println("Compare [", strings.TrimSpace(strings.ToUpper(allowExtenstion)), "] with [", strings.TrimSpace(strings.ToUpper(fileExtension)), "]")
		if strings.TrimSpace(strings.ToUpper(allowExtenstion)) == strings.TrimSpace(strings.ToUpper(fileExtension)) {
			fmt.Println("Allowed extention ")
			return true
		}
	}
	fmt.Println("Not Allowed extention ==>", fileExtension)
	return false

}

func saveAttachment(filename string, value []byte, index int) string {

	filePath := generatePathAttachment(index)
	filePathAttachment := filepath.Join(filePath, filename)
	fmt.Println("File path 3", filePath)

	buf := new(strings.Builder)
	io.Copy(buf, bytes.NewReader(value))
	if err := ioutil.WriteFile(filePathAttachment, []byte(buf.String()), 0644); err != nil {
		fmt.Println("ERROR Save Attachment == ", err.Error())
	}

	return filePath
}

func generatePathAttachment(index int) string {

	destination := "attachment"
	path, _ := os.Getwd()

	filePath := filepath.Join(path, destination)
	fmt.Println("File path 1", filePath)
	os.Mkdir(filePath, 0755)

	filePath = filepath.Join(filePath, strconv.Itoa(index))
	fmt.Println("File path 2", filePath)
	os.Mkdir(filePath, 0755)

	return filePath
}

func saveToDatabase(emailModel repository.TicketEmail, attachModel repository.TicketEmailAttachment) {
	var err errorModel.ErrorModel
	var idTicket int64
	db := serverConfig.ServerAttribute.DBConnection

	emailModel.CreatedAt.Time = time.Now()
	emailModel.CreatedBy.Int64 = 0
	emailModel.UpdatedAt.Time = time.Now()
	emailModel.UpdateBy.Int64 = 0
	emailModel.Deleted.Bool = false

	fmt.Println("================== SAVE ISI DB ========================")
	fmt.Println("-> ", emailModel)
	fmt.Println("-> ", attachModel)

	err, idTicket = dao.EmailDAO.InsertEmail(db, emailModel)

	if err.Error != nil {
		fmt.Println("Error => ", err)
		return
	}

	err = dao.AttachDAO.InsertAttachment(db, attachModel, idTicket)

	fmt.Println("Save success  ID ==>", idTicket)
	// return

}
