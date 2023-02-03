package mail

import (
	"os"

	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

type Mail struct {
	To          string `json:"to"`
	CC          string `json:"cc"`
	BCC         string `json:"bcc"`
	Subject     string `json:"subject"`
	Body        string `json:"body"`
	HTMLBody    string `json:"htmlbody"`
	Attachments string `json:"attachments"`
}

func CreateMail(mailObj Mail) *ole.IDispatch {
	ole.CoInitialize(0)
	unknown, _ := oleutil.CreateObject("Outlook.Application")
	outlook, _ := unknown.QueryInterface(ole.IID_IDispatch)
	ns := oleutil.MustCallMethod(outlook, "GetNamespace", "MAPI").ToIDispatch()
	oleutil.MustCallMethod(ns, "Logon").ToIDispatch()

	mail := oleutil.MustCallMethod(outlook, "CreateItem", 0).ToIDispatch()
	if mailObj.To != "" {
		oleutil.PutProperty(mail, "To", mailObj.To)
	}
	if mailObj.CC != "" {
		oleutil.PutProperty(mail, "CC", mailObj.CC)
	}
	if mailObj.BCC != "" {
		oleutil.PutProperty(mail, "BCC", mailObj.BCC)
	}
	if mailObj.Subject != "" {
		oleutil.PutProperty(mail, "Subject", mailObj.Subject)
	}
	if mailObj.Body != "" {
		oleutil.PutProperty(mail, "Body", mailObj.Body)
	}
	if mailObj.HTMLBody != "" {
		oleutil.PutProperty(mail, "HTMLBody", mailObj.HTMLBody)
	}
	if mailObj.Attachments != "" {
		attachments := oleutil.MustGetProperty(mail, "Attachments").ToIDispatch()
		cwd, _ := os.Getwd()
		oleutil.MustCallMethod(attachments, "Add", cwd+"\\"+mailObj.Attachments).ToIDispatch()
	}

	return mail
}

func SendMail(mail *ole.IDispatch) {
	oleutil.MustCallMethod(mail, "Send")
}
