// Processing

IF TEST = 1
{
	DYNAMIC i
	i.friendlyName = ""
	i.url = ""

	ASSIGN token = "Bearer adc3cc3efd434fbdc4980bbe1ba4fdba6f0744c3"
	ASSIGN baseURL = "https://api-de-NA1.niceincontact.com"
	
	ASSIGN postId = "email_ecolabtest1-eventusg-abc123@mail-de-NA1.niceincontact.com_9f41e35d-fa64-46be-9ff1-b50da953f686"
	ASSIGN caseId = "127035285501"
}

dfoerror = false
errormsg = ""

// This user is used to show "who" made the changes in DFO
// Do not inactivate this user without changing the adminuser
adminuser = "13650" // Jeff Taylor

// This snippet uses the DFO "Reply" to forward an e-mail to Conexiom
IF Conexiom
{
	proxy = GetRESTProxy()
	proxy.ContentType = "application/json"
	proxy.AddHeader("Authorization", "{token}")

	// this is the Conexiom e-mail address that we forward e-mails to
	//forwardto = "test-ecolabinc.us@conexiom.net"
	forwardto = "EcolabInc.us@conexiom.net"

	messageId = "{postId}"

	// deal with quotes, carriage returns, line feeds, and add Case# to the mail body
	mailbody = "{mailbody.replace('"','\"')}"
	mailbody = "{mailbody.replace(char(13),'<br>')}"
	mailbody = "{mailbody.replace(char(10),'')}"
	mailbody = "[Case #{caseId}]<br>{mailbody}"

	// \"attachments\": [ { \"friendlyName\": \"kitten.png\", \"url\": \"https://placekitten.com/200/300\" } ]"

	c = 0
	ASSIGN att = $", \"attachments\": [ "
	FOREACH i IN mailattach
	{
		c = c + 1
		IF c > 1
		{
			att = "{att}, "
		}
		att = $"{att} \{ \"friendlyName\": \"{i.friendlyName}\", \"url\": \"{i.url}\" } "
	}
	att = $"{att} ] "

	// blank out string if no attachments
	IF c = 0
	{
		att = ""
	}

	jsonbody = $"\{ \"channelId\": \"{channelId}\", \"replyToMessageId\": \"{messageId}\", \"authorId\": {adminuser}, \"endUserRecipients\": [ \{ \"idOnExternalPlatform\": \"{forwardto}\", \"name\": \"John Doe\", \"isPrimary\": true, \"isPrivate\": false } ], \"isForwarded\": true, \"title\": \"{mailsubject}\", \"messageType\": \"TEXT\", \"payload\": \{ \"text\": \"{mailbody}\" } {att} }"

	ASSIGN RequestParms = ""
	ASSIGN URI="{baseURL}/engager/2.0/posts/{postId}/reply"

	postResult = proxy.MakeRestRequest(URI, jsonbody, 0, "POST")

	ASSIGN HTTPStatus="{Proxy.StatusCode}"

	SELECT
	{
		CASE HTTPStatus = 201 
		{
			// ok
		}

	//	CASE httpStatus = 401 //Invalid session
	//	CASE HTTPStatus = 500 // internal error
		
		DEFAULT
		{
			dfoerror = true
			errormsg = "{errormsg}<br><br>There was an error forwarding the e-mail to Conexiom for CaseID {caseId} (HTTPSTATUS={httpstatus}).  You will only receive this message once per hour."
		}
	}
}

// This snippet uses the DFO "Reply" to send an acknowledgement e-mail to customer
IF sendack
{
	proxy = GetRESTProxy()
	proxy.ContentType = "application/json"
	proxy.AddHeader("Authorization", "{token}")

	forwardto = "{mailfrom}"
	messageId = "{postId}"

	// deal with quotes, carriage returns, line feeds, and add Case# to the mail body
//	mailbody = "{mailbody.replace('"','\"')}"
//	mailbody = "{mailbody.replace(char(13),'<br>')}"
//	mailbody = "{mailbody.replace(char(10),'')}"
	mailbody = "{ackmessage}"
	
	jsonbody = $"\{ \"channelId\": \"{channelId}\", \"replyToMessageId\": \"{messageId}\", \"authorId\": {adminuser}, \"endUserRecipients\": [ \{ \"idOnExternalPlatform\": \"{forwardto}\", \"name\": \"John Doe\", \"isPrimary\": true, \"isPrivate\": false } ], \"isForwarded\": false, \"title\": \"{mailsubject}\", \"messageType\": \"TEXT\", \"payload\": \{ \"text\": \"{mailbody}\" } }"

	ASSIGN RequestParms = ""
	ASSIGN URI="{baseURL}/engager/2.0/posts/{postId}/reply"

	postResult = proxy.MakeRestRequest(URI, jsonbody, 0, "POST")

	ASSIGN HTTPStatus="{Proxy.StatusCode}"

	SELECT
	{
		CASE HTTPStatus = 201
		{
			// ok
		}

	//	CASE httpStatus = 401 //Invalid session
	//	CASE HTTPStatus = 500 // internal error
		
		DEFAULT
		{
			dfoerror = true
			errormsg = "{errormsg}<br><br>There was an error sending the acknowledgment for CaseID {caseId} (HTTPSTATUS={httpstatus}).  You will only receive this message once per hour."
		}
	}
}



// Assign to Queue
IF changequeue
{
	proxy2 = GetRESTProxy()
	proxy2.ContentType = "application/json"
	proxy2.AddHeader("Authorization", "{token}")

	ASSIGN RequestParms = $"\{ \"routingQueueId\": \"{targetqueue}\", \"assignedBy\": {adminuser} }"

	ASSIGN URI="{baseURL}/engager/2.0/posts/{postId}/cases/{caseId}/routing-queue-assignment"

	putQueueResult = proxy2.MakeRestRequest(URI, requestParms, 0, "PUT")

	ASSIGN HTTPStatus="{Proxy2.StatusCode}"

	SELECT
	{
		CASE HTTPStatus = 200 // ok
		{
			// Success
		}

	//	CASE httpStatus = 401 //Invalid session

		DEFAULT
		{
			dfoerror = true
			errormsg = "{errormsg}<br><br>There was an error assigning the queue for CaseID {caseId} (HTTPSTATUS={httpstatus}).  You will only receive this message once per hour."
		}
	}
}


// Assign to Specific User
IF changeuser
{
	proxy3 = GetRESTProxy()
	proxy3.ContentType = "application/json"
	proxy3.AddHeader("Authorization", "{token}")

	ASSIGN RequestParms = $"\{ \"userId\": {targetuser}, \"assignedBy\": {adminuser} }"

	ASSIGN URI="{baseURL}/engager/2.0/posts/{postId}/cases/{caseId}/inbox-assignment"

	putUserResult = proxy3.MakeRestRequest(URI, requestParms, 0, "PUT")

	ASSIGN HTTPStatus="{Proxy3.StatusCode}"

	SELECT
	{
		CASE HTTPStatus = 200 // ok
		{
			// Success
		}

	//	CASE httpStatus = 401 //Invalid session
	//	{

		DEFAULT
		{
			dfoerror = true
			errormsg = "{errormsg}<br><br>There was an error assigning the dedicated preferred user for CaseID {caseId} (HTTPSTATUS={httpstatus}).  You will only receive this message once per hour."
		}
	}
}

// Set Priority
IF priority > 0
{
	proxy4 = GetRESTProxy()
	proxy4.ContentType = "application/json"
	proxy4.AddHeader("Authorization", "{token}")

	ASSIGN RequestParms = $"\{ \"routingQueuePriority\": {priority} }"

	ASSIGN URI="{baseURL}/dfo/3.0/contacts/{caseId}"

	putPriorityResult = proxy4.MakeRestRequest(URI, requestParms, 0, "PUT")

	ASSIGN HTTPStatus="{Proxy4.StatusCode}"

	SELECT
	{
		CASE HTTPStatus = 200 // ok
		CASE HTTPStatus = 204 // ok
		{
			// Success
		}

		DEFAULT
		{
			dfoerror = true
			errormsg = "{errormsg}<br><br>There was an error setting the priority for CaseID {caseId} (HTTPSTATUS={httpstatus}).  You will only receive this message once per hour."
		}
	}
}

// Update Status
IF changestatus
{
	proxy5 = GetRESTProxy()
	proxy5.ContentType = "application/json"
	proxy5.AddHeader("Authorization", "{token}")

	ASSIGN RequestParms = $"\{ \"status\": \"{targetstatus}\" }"

	ASSIGN URI="{baseURL}/engager/2.0/posts/{postId}/cases/{caseId}/status"

	putStatusResult = proxy5.MakeRestRequest(URI, requestParms, 0, "PUT")

	ASSIGN HTTPStatus="{Proxy5.StatusCode}"

	SELECT
	{
		CASE HTTPStatus = 200 // ok
		CASE HTTPStatus = 204 // ok
		{
			// Success
		}

		DEFAULT
		{
			dfoerror = true
			errormsg = "{errormsg}<br><br>There was an error updating the status for CaseID {caseId} (HTTPSTATUS={httpstatus}).  You will only receive this message once per hour."
		}
	}
}







// Pest Rules

// Get Case Info has assigned the following variables
//
// mailfrom
// mailsubject
// mailbody

// Case Status:
//	new
//	open
//	pending
//	escalated
//	resolved
//	closed
//	trashed

IF test
{
//	FUNCTION BeginsWith(s1, s2)
//	{
//		l = s2.length()
//		s3 = s1.left(l)
//			
//		IF s3.contains(s2)
//		{
//			RETURN true
//		} 
//		ELSE
//		{
//			RETURN false
//		}
//	}
//	mailfrom = "Noreply@officetrax.com"
//	mailto = "pest@ecolab.com"
//  pestcscsmiles@ecolab.com
//	mailsubject = "Work Order Close out for WKO-01599718 Store 29215 with Trade of Pest Control"
//	mailbody = "something else here"
}

ASSIGN changequeue = false
ASSIGN changeuser = false
ASSIGN changestatus = false
ASSIGN priority = 0
ASSIGN sendack = false
ASSIGN debug = 0

ASSIGN cont = true

ASSIGN Conexiom = false

SELECT
{
	// Auto Complete Rules 
	CASE mailfrom.contains("MicrosoftExchange")
	CASE mailsubject.contains("Auto answer")
	CASE mailsubject.contains("Auto Reply")
	CASE mailsubject.contains("Auto response")
	CASE mailsubject.contains("auto-reply")
	CASE mailfrom.contains("Auto-Reply")
	CASE mailfrom.contains("auto-sender")
	CASE mailsubject.contains("Automated response")
	CASE mailsubject.contains("Automatic Omnigate Message")
	CASE mailsubject.contains("Automatic reply")
	CASE mailsubject.contains("Automatic response")
	CASE mailsubject.contains("Automaticka odpoved")
	CASE mailsubject.contains("AUTOMATICKA ODPOVID")
	CASE mailsubject.contains("Automatisch antwoord")
	CASE mailsubject.contains("Automatisk_svar")
	CASE mailsubject.contains("Automatsvar")
	CASE mailsubject.contains("AutoReply")
	CASE mailsubject.contains("AutoResp")
	CASE mailfrom.contains("autoresponder")
	CASE mailsubject.contains("Autoresponse")
	CASE mailsubject.contains("Autosvar")
	CASE mailsubject.contains("away from my mail")
	CASE mailsubject.contains("away from the office")
	CASE mailsubject.contains("Delivery notification")
	CASE mailsubject.contains("Details of my business trips")
	CASE mailsubject.contains("E-mail Received!")
	CASE mailsubject.contains("Ecolab Spam Filter")
	CASE mailsubject.contains("Keep more of what you make!")
	CASE mailsubject.contains("Your message was received")
	CASE mailfrom.contains("omsnotification@itradenetwork.com") & mailsubject.contains("ITN Inbound 810 Transaction Error")
	CASE mailfrom.contains("IFM-Invoices@eu.jll.com") & mailsubject.contains("Acknowledgement: Ecolab Invoice")
	CASE mailfrom.contains("PESTWOMGMT@ecolab.com") & mailsubject.contains("PPM work Order Generated")
	CASE mailfrom.contains("Mediclean.Invoices@uk.issworld.com") & mailsubject.contains("Re: PLEASE READ: IMPORTANT INFORMATION REDARDING YOUR EMAIL")
	CASE mailbody.left(32) = "Thank you for contacting Ecolab."
	CASE mailbody.left(28) = "This is an automated message"
	CASE mailfrom.contains("accountpayables@clevelandcliffs.com") & mailsubject.contains("Thank you for submitting your invoice electronically")
//	CASE mailsubject.contains("Thank you for Submitting a Case to REPAY")
//	CASE mailsubject.contains("#00193126")
//	CASE mailfrom.contains("dfomailtest14@ecolab.com")
//	CASE mailbody.contains("This is an automated message confirming receipt of your request, please do not reply to this message.")
//	CASE mailfrom.contains("dfomailtest14@ecolab.com")
//	CASE mailsubject.contains("New Process FOR Submitting Comp Requests - Please read")
//	
	
	{
		changequeue = true
		targetqueue = "a82590c9-41c9-42b7-9a45-edca9bb87c37"	// CS.US.PEST.EN.CB.EM.WorkOrders
		
		changestatus = true
		targetstatus = "closed"
		
		cont = false
		debug = 1
	}

	// DEQ Rules
	CASE mailfrom.contains("-MaiSer-")
	CASE mailsubject.contains("Abwesenheitsnotiz")
	CASE mailsubject.contains("Addressee Unavailable")
	CASE mailsubject.contains("Adressänderung")
	CASE mailsubject.contains("bad-style address")
	CASE mailfrom.contains("badaddress")
	CASE mailsubject.contains("Bevestiging ontvangen")
	CASE mailsubject.contains("bounced message")
	CASE mailfrom.contains("ccmail_agent")
	CASE mailsubject.contains("Conversion fail")
	CASE mailsubject.contains("could not send message")
	CASE mailsubject.contains("Delivery Error")
	CASE mailsubject.contains("Delivery Failed")
	CASE mailsubject.contains("Delivery Failure")
	CASE mailsubject.contains("Delivery Problem Notification")
	CASE mailsubject.contains("Delivery Returned")
	CASE mailsubject.contains("Delivery Status Notification")
	CASE mailsubject.contains("Delivery-Report")
	CASE mailsubject.contains("Dìkuji za maila")
	CASE mailsubject.contains("E-mail Unavailable")
	CASE mailsubject.contains("Error Response")
	CASE mailsubject.contains("Error sending mail")
	CASE mailsubject.contains("Extended Absence Response")
	CASE mailsubject.contains("Failed mail")
	CASE mailsubject.contains("failed message delivery")
	CASE mailsubject.contains("failure notice")
	CASE mailsubject.contains("Ihre Mail")
	CASE mailsubject.contains("Inaccessible e-mail address")
	CASE mailsubject.contains("INBOUND MESSAGE ERR")
	CASE mailsubject.contains("Invalid mailbox")
	CASE mailsubject.contains("Invalid user")
	CASE mailsubject.contains("Mail Did Not Get Through")
	CASE mailsubject.contains("Mail error")
	CASE mailsubject.contains("Mail failed")
	CASE mailsubject.contains("Mail failure")
	CASE mailsubject.contains("Mail recipient has left Enter-Net")
	CASE mailfrom.contains("Mail-Gateway")
	CASE mailfrom.contains("mailer-daemon")
	CASE mailfrom.contains("mail_master")
	CASE mailfrom.contains("mdaemon")
	CASE mailsubject.contains("message failed")
	CASE mailsubject.contains("message not sent")
	CASE mailsubject.contains("message rejected")
	CASE mailsubject.contains("message was not sent")
	CASE mailsubject.contains("NDN:")
	CASE mailsubject.contains("no interest !!")
	CASE mailsubject.contains("No such user")
	CASE mailsubject.contains("Non-Delivery")
	CASE mailsubject.contains("Non-existing employee")
	CASE mailsubject.contains("Nondeliverable")
	CASE mailsubject.contains("Not a WORLDPATH client")
	CASE mailsubject.contains("not deliverable")
	CASE mailsubject.contains("not delivered")
	CASE mailsubject.contains("not_a_jono_addy")
	CASE mailsubject.contains("Odpoved na zpravu")
	CASE mailsubject.contains("Ontvangstbevestiging")
	CASE mailsubject.contains("ponse_automatique")
	CASE mailfrom.contains("postadm")
	CASE mailfrom.contains("postmast")
	CASE mailfrom.contains("postmaster")
	CASE mailsubject.contains("problem delivering your mail")
	CASE mailsubject.contains("Response from Administrator")
	CASE mailsubject.contains("Response from bdbad")
	CASE mailsubject.contains("Response from rlozano")
	CASE mailsubject.contains("Resposta Automatica")
	CASE mailsubject.contains("Return message")
	CASE mailsubject.contains("Returned Mail")
	CASE mailsubject.contains("Returned mail: see transcript for details")
	CASE mailsubject.contains("Returned to Sender")
	CASE mailsubject.contains("Réponse automatique")
	CASE mailsubject.contains("Service Message")
	CASE mailsubject.contains("SMS error response")
	CASE mailsubject.contains("SMS message")
	CASE mailfrom.contains("supervisor")
	CASE mailsubject.contains("Thanks for writing ER!")
	CASE mailsubject.contains("Thanks for your e-mail message!")
	CASE mailsubject.contains("Troubles delivering the message")
	CASE mailsubject.contains("Unable to deliver mail")
	CASE mailfrom.contains("unknown")
	CASE mailsubject.contains("unknown address")
	CASE mailsubject.contains("unknown domain")
	CASE mailsubject.contains("unknown recipient")
	CASE mailsubject.contains("User Not at VISTA.COM Domain")
	CASE mailsubject.contains("user not found")
	CASE mailsubject.contains("user unknown")
	CASE mailsubject.contains("Warning - delayed mail")
	CASE mailsubject.contains("X.400 Inter-Personal Notification")
	CASE mailsubject.contains("Your Message To Juno")
	CASE mailsubject.contains("ZAZ Reply")
	CASE mailsubject.contains("#11264")
	
	{
		changequeue = true
		targetqueue = "1124e64f-759c-48da-a7f0-f38ad6c31c5c"	// DEQ.PEST.SS

		cont = false
		debug = 2
	}
}

IF cont
{
	// Auto Ack
	SELECT
	{
		// Accept
//		CASE mailfrom.contains("ServiceRequest@scalert.com")
//		{
//			sendack = true
//			ackmessage = "Accept<br><br>"
//			
//			changestatus = true
//			targetstatus = "new"
//			debug = 3
//		}
		
		// Auto Ack rules are reversed from ECE logic
		CASE mailfrom.contains("fservice@starbucks.com")
		CASE mailsubject.contains("Canada Pest Message")
		//CASE mailfrom.contains("ServiceRequest@scalert.com")
		CASE mailfrom.contains("cs-specialtyorders@ecolab.com")
		CASE mailfrom.contains("CS-Workorders@ecolab.com")
		CASE mailfrom.contains("CAN_Product_Order_Noreply@ecolab.com")
		CASE mailfrom.contains("information@verisae.com")
		CASE mailfrom.contains("ECOSERV@ecolab.com")
		CASE mailfrom.contains("facilitiesHELP@brookshires.com")
		CASE mailfrom.contains("CS-InstitutionalOrders@ecolab.com")
		CASE mailfrom.contains("swisher.csdistributors@ecolab.com")
		CASE mailfrom.contains("swisher.csorders@ecolab.com")
		CASE mailfrom.contains("CS-INSTOrderForm@ecolab.com")
		CASE mailfrom.contains("noreply@retarus.net")
		CASE mailfrom.contains("problemticket@medline.com")
		CASE mailfrom.contains("metatradeautoreply@ghx.com")
		CASE mailfrom.contains("NoReply@pfgc.com")
		CASE mailfrom.contains("noreply@hcparts.ecolab.com")
		CASE mailfrom.contains("Reply")
		CASE mailfrom.contains("WholeFoods@verisae.com")
		CASE mailfrom.contains("workflowapproval@workflow.innout.com")
		CASE mailfrom.contains("gei-prod-sender@imsevolve.com")
		{
			// Do Nothing
		}
		
		DEFAULT
		{
			sendack = true
			ackmessage = "Thank you for contacting Ecolab.<br><br>This is an automated message confirming receipt of your request, please do not reply to this message.  We will begin processing your request as soon as possible. If your request is urgent please contact Customer Service at 1-800-325-1671.<br><br>Thank you,<br>"
			
			changestatus = true
			targetstatus = "new"
		}
	}
	
	// Pest Routing
	SELECT
	{
		// Auto Close
		CASE mailsubject.beginswith("EditWO")
		CASE mailsubject.beginswith("Jack In The Box/Qdoba has authorized")
		CASE mailsubject.beginswith("WO Note / Wal-Mart Stores")
		CASE mailsubject.contains("Completed: Chuck E. Cheese Entertainment")
		CASE mailsubject.contains("Scheduled Maintenance plan: Chuck E. Cheese Entertainment")
		CASE mailsubject.contains("Signed Off: Chuck E. Cheese Entertainment")
		CASE mailsubject.contains("PPM Work Order Generated")
		CASE mailsubject.beginswith("Jack In The Box/Qdoba has authorized the invoice for Service Request")
		CASE mailsubject.beginswith("StarBoard-Wendy&amp;amp;amp;amp;amp;amp;amp;amp;#039;s has authorized the invoice for Service Request")
		CASE mailsubject.beginswith("KBP FOODS has authorized the invoice for Service Request")
		CASE mailsubject.beginswith("Flagged Service Request Alert")
		CASE mailsubject.contains("Scheduled/Project (4 day onsite/2 week resolution)")
		CASE mailsubject.contains("Jack In The Box/Qdoba has sent you a message concerning Service Request")
		CASE mailsubject.beginswith("JLL BNSF has authorized the invoice for Service Request")
		CASE mailsubject.contains("Signed Off: Kum &amp;amp;amp;amp;amp;amp; Go")
		CASE mailsubject.contains("New Notes added for location")
		CASE mailsubject.contains("Subject: Jack In The Box/Qdoba has approved your quote for Service Request")
		CASE mailsubject.beginswith("Survey Invitation")
		CASE mailfrom.contains("ePASS@supervalu.com")
		CASE mailsubject.beginswith("Jack In The Box/Qdoba has changed the NTE for Service Request")
		CASE mailsubject.beginswith("JLL BNSF has sent you a message concerning Service Request")
		CASE mailsubject.contains("has been Auto-updated to a status of Approved: Kum &amp;amp;amp;amp;amp;amp; Go L.C.: Kum and Go L.C")
		CASE mailfrom.contains("centraldisburse@supervalu.com")
		CASE mailsubject.beginswith("CBRE Sprint Work Order Network has changed the due by date/time for Service Request")
		CASE mailfrom.contains("janm@unitedagcoop.com")
		CASE mailfrom.contains("kp-ap-customer@kp.org")
		CASE mailsubject.contains("update processing not completed")
		CASE mailsubject.beginswith("H-E-B My Facility Preventative Work order")
		CASE mailsubject.contains("Thank you for contacting the Accounts Payable department for Stanford Health Care")
		CASE mailsubject.contains("WO Note Planned Maintenance / Islands Restaurants / 014 - Carmel Mountain / CA / ECOLAB PEST")
		CASE mailsubject.contains("Your Email Has Been Received")
		CASE mailsubject.contains("Two folders in your mailbox have the same name")
		CASE mailsubject.contains("There is a new comment on AskMe request")
		CASE mailsubject.contains("Please Read This Message in Its Entirety")
		CASE mailsubject.contains("Case AP0244471 Is Resolved")
		CASE mailbody.contains("I will be out of the office")
		CASE mailsubject.contains("Custom reply, please do not respond.")
		CASE mailsubject.contains("Your ticket has been created")
		CASE mailsubject.contains("You have been assigned a new job or jobs with a badge")
		CASE mailsubject.contains("Problem Description: Preventative Maintenance - Pest Control (Monthly)")
		CASE mailsubject.contains("Fwd: Service Request - [#XN12824]")
		CASE mailfrom.contains("usrcs@ecolab.com")
		CASE mailsubject.contains("Approved to Pay by MGM Resorts International")
		CASE mailsubject.contains("MGM Resorts International Revised Purchase Order")
		CASE mailsubject.contains("Two folders in your mailbox have the same name")
		CASE mailfrom.contains("AUS-APSS@aramark.com")
		CASE mailfrom.contains("apinquiries@cecentertainment.com")
		CASE mailsubject.contains("ECOLAB Service Report/Invoice 0099184 PANERA BREAD")
		CASE mailsubject.contains("Service Report/Invoice 0098795")
		CASE mailsubject.contains("An open-item stateme Delivery Confirmation Delivery Confirmation Delivery Confirmation Delivery Confirmation Delivery Confirmation Del")
		CASE mailsubject.contains("Your ticket has been Closed")
		CASE mailsubject.contains("Pending Quote Revisions")
		CASE mailfrom.contains("MicrosoftExchange")
		CASE mailfrom.contains("jpmorganchase@ilrs.360facility.net") & mailbody.contains("UPDATE: Valid CD extension found.")
		CASE mailfrom.contains("dnap@delawarenorth.com") & mailbody.contains("Your ticket has been Closed")
		CASE mailfrom.contains("dollartree.admin@officetrax.com") & mailbody.contains("Request Close Out")
		CASE mailfrom.contains("information@verisae.com") & mailbody.contains("Estimate has been Approved")
		CASE mailfrom.contains("requestfmsupport@heb.com") & mailbody.contains("Preventative Issue type and priority of Preventative")
		CASE mailfrom.contains("gei-prod-sender@ims-evolve.com") & mailbody.contains("Priority: PPM")
		CASE mailfrom.contains("mail-robot@workoasis.net") & mailbody.contains("has been completed")
		CASE mailfrom.contains("mail-robot@workoasis.net") & mailsubject.contains("has been created from a Scheduled Maintenance")
		CASE mailfrom.contains("DoNotReply@conduent.com") & mailsubject.contains("Thank you for submitting your invoice")
		CASE mailfrom.contains("provideralert@workordernetwork.corrigo.com") & mailsubject.contains("CBRE Sprint Work Order Network has authorized")
		CASE mailfrom.contains("gei-prod-sender@ims-evolve.com") & mailsubject.contains("Invoice Required For Work Order")
		CASE mailfrom.contains("brookdale@ilrs.360facility.net") & mailsubject.contains("SCHEDULED Brookdale | Insects/Bugs")
		CASE mailfrom.contains("compass@starbucks.com") & mailsubject.contains("SUBMISSION ERROR:")
		CASE mailfrom.contains("snprod@medtronic.com") & mailsubject.contains("Comments Added")
		CASE mailfrom.contains("officetrax.noreply@corp.core7.com")
		CASE mailfrom.contains("brookdale@ilrs.360facility.net") & mailsubject.contains("UPDATE")
		CASE mailfrom.contains("WO Note") & mailsubject.contains("ServiceRequest@scalert.com")
		CASE ( mailto.contains("CSWorkorders@ecolab.com") | mailcc.contains("CSWorkorders@ecolab.com") ) & mailfrom.contains("CSWorkorders@ecolab.com") & mailbody.contains("This is an automated message confirming receipt of your request")
		CASE ( mailto.contains("CSC.Workorders@ecolab.com") | mailcc.contains("CSC.Workorders@ecolab.com") ) & mailfrom.contains("CSWorkorders@ecolab.com") & mailbody.contains("This is an automated message confirming receipt of your request")
		CASE ( mailto.contains("CSCSupervisors@ecolab.com") | mailcc.contains("CSCSupervisors@ecolab.com") ) & mailfrom.contains("CSWorkorders@ecolab.com") & mailbody.contains("This is an automated message confirming receipt of your request")
		CASE ( mailto.contains("CSWorkorders@ecolab.com") | mailcc.contains("CSCWorkorders@ecolab.com") ) & mailfrom.contains("CSCWorkorders@ecolab.com") & mailbody.contains("This is an automated message confirming receipt of your request")
		CASE ( mailto.contains("CSWorkorders@ecolab.com") | mailcc.contains("CS-Workorders@ecolab.com") ) & mailfrom.contains("CS-Workorders@ecolab.com") & mailbody.contains("This is an automated message confirming receipt of your request")
		CASE ( mailto.contains("pestcscsmiles@ecolab.com") | mailcc.contains("pestcscsmiles@ecolab.com") ) & mailfrom.contains("CS-Workorders@ecolab.com") & mailbody.contains("This is an automated message confirming receipt of your request")
		CASE ( mailto.contains("pestcscsmiles@ecolab.com") | mailcc.contains("pestcscsmiles@ecolab.com") ) & mailfrom.contains("CS-Workorders@ecolab.com") & mailbody.contains("This is an automated message confirming receipt of your request")
		CASE mailfrom.contains("mail-robot@workoasis.net") & mailsubject.contains("has been Signed Off")
		CASE mailfrom.contains("mail-robot@workoasis.net") & mailsubject.contains("Cancelled")
		CASE mailfrom.contains("Completed – waiting sign off") & mailsubject.contains("mail-robot@workoasis.net")
		CASE mailfrom.contains("do_not_reply@mgmresorts.coupahost.com") & mailsubject.contains("New PO")
		CASE mailfrom.contains("medtronicprod@service-now.com") & mailsubject.contains("Comments Added")
		CASE mailfrom.contains("APInvoice@vailresorts.com") & mailsubject.contains("Your Submission Has Been Received")
		CASE mailfrom.contains("ServiceRequest@scalert.com") & mailbody.contains("Pest Control: Inclusive Package (Auto Complete)")
		CASE mailfrom.contains("APInvoice@vailresorts.com") & mailsubject.contains("Your Submission Has Been Received")
		CASE mailfrom.contains("operations@managepathcloud.com") & mailsubject.contains("Completed Work")
		CASE mailfrom.contains("apinquiries@boydgaming.com") & mailsubject.contains("**Please Read**")
		CASE mailfrom.contains("request@ecotrak.com") & mailsubject.contains("Update Proposal")
		CASE mailfrom.contains("AUS-APSS@aramark.com") & mailsubject.contains("Thank you for your inquiry")
		CASE mailfrom.contains("ServiceRequest@scalert.com") & mailbody.contains("Billing Purpose Only")
		CASE mailfrom.contains("accountspayable@greatfoodserv.com") & mailsubject.contains("Your Invoice Has Been Received")
		CASE mailfrom.contains("nvinvoices@boydgaming.com") & mailsubject.contains("Please Read")
		CASE mailfrom.contains("PESTWOMGMT@ecolab.com") & mailsubject.contains("Rodents and Insects assigned to Kannha Evans for Completed Work")
		CASE mailfrom.contains("operations@managepathcloud.com") & mailsubject.contains("Dispatched")
		CASE mailfrom.contains("dispatch@officetrax.com") & mailbody.contains("Cockroach/Rodent Program/One-Shot Service")
		CASE mailfrom.contains("officetrax.noreply@corp.core7.com") & mailsubject.contains("Work Order Close out")
		CASE mailfrom.contains("supportdesk@sprouts.com") & mailsubject.contains("Comments on a Closed tickets are Not Monitored")
		CASE mailfrom.contains("do_not_reply@mgmresorts.coupahost.com") & mailsubject.contains("marked as paid")
		CASE mailfrom.contains("ecoservereporting@ecolab.com") & mailsubject.contains("List of Messages on Hold and need to be forwarded")
		CASE mailfrom.contains("usrcs@ecolab.com") & ( mailto.contains("CSCSupervisor@ecolab.com") | mailcc.contains("CSCSupervisor@ecolab.com") ) & mailsubject.contains("CSCSupervisor@ecolab.com")
		CASE mailfrom.contains("mail-robot@workoasis.net") & mailsubject.contains("has been approved")
		CASE mailfrom.contains("information@tradeshift.com") & mailsubject.contains("CBRE on Behalf of APM ask that you accept purchase order")
		CASE mailfrom.contains("operations@managepathcloud.com") & mailsubject.contains("Pending Quote Revisions")
		CASE mailsubject.contains("for Dispatched")
		CASE mailfrom.contains("provideralert@workordernetwork.corrigo.com") & mailsubject.contains("JLL BNSF has approved your quote for Service Request")
		CASE mailfrom.contains("dispatch@officetrax.com") & mailsubject.contains("FOR BILLNG PURPOSES ONLY")
		CASE mailfrom.contains("provideralert@workordernetwork.corrigo.com") & mailsubject.contains("StarBoard-Wendy&amp;amp;amp;#039;s has approved your quote for Service Request")
		CASE mailfrom.contains("provideralert@workordernetwork.corrigo.com") & mailbody.contains("| MONTHLY |")
		CASE mailsubject.contains("ECOLAB Service Report/Invoice") & mailbody.contains("THIS IS A SYSTEM AUTO-REPLY MESSAGE")
		CASE mailfrom.contains("ServiceRequest@scalert.com") & mailsubject.contains("WO cancelled")
		CASE mailbody.contains("Description1: ***For billing purposes only")
		CASE mailfrom.contains("ServiceRequest@scalert.com") & mailbody.contains("Description1: ***For billing purposes only")
		CASE mailbody.contains("Description1: ***For billing purposes only") & mailsubject.contains("Email address update")
		CASE mailfrom.contains("PESTWOMGMT@ecolab.com") & mailsubject.contains("assigned to Kannha Evans")
		CASE mailfrom.contains("supportdesk@sprouts.com") & mailsubject.contains("ETA has changed")
		CASE mailfrom.contains("support@accuinc.com") & mailsubject.contains("General Inquiry From Email")
		CASE mailfrom.contains("help@chainrs.com") & mailsubject.contains("CRS Needs ETA For PM Dispatch")
		CASE mailfrom.contains("help@chainrs.com") & mailsubject.contains("CRS Needs an ETA for Work Requested")
		CASE mailfrom.contains("request+noreply@ecotrak.com") & mailsubject.contains("Invoice Update: California Fish Grill")
		CASE mailfrom.contains("supportdesk@sprouts.com") & mailsubject.contains("ready to bill")
		CASE mailfrom.contains("Pest-AuthCredits@ecolab.com") & mailsubject.contains("Credit Request has been received")
		CASE mailfrom.contains("supportdesk@sprouts.com") & mailsubject.contains("ETA has changed")
		CASE mailfrom.contains("supportdesk@sprouts.com") & mailsubject.contains("rejected")
		CASE mailfrom.contains("support@accuinc.com") & mailsubject.contains("RE: Service Request")
		CASE mailfrom.contains("request@ecotrak.com") & mailsubject.contains("PEST CONTROL MONTHLY SERVICE")
		CASE mailfrom.contains("request+noreply@ecotrak.com") & mailsubject.contains("Re: Re: New Service")
		CASE mailfrom.contains("compass@starbucks.com") & mailsubject.contains("SUBMISSION ERROR")
		CASE mailfrom.contains("Alert@scalert.com") & mailsubject.contains("Unsatisfactory Feedback")
		CASE mailfrom.contains("apinvoices@brooksbrothers.com") & mailsubject.contains("Confirmation Email - AP Invoices")
		CASE mailfrom.contains("CDS.Payables@cirquedusoleil.com") & mailsubject.contains("RE:ECOLAB Service Report/Invoice")
		CASE mailfrom.contains("Payables.SSI.IHQ@cirquedusoleil.com") & mailbody.contains("Merci de nous avoir écrit. Nous avons bien reçu votre message et nous communiquerons avec vous dès que possible")
		CASE mailfrom.contains("Noreply@officetrax.com") & mailsubject.contains("Work Order Close")
		CASE mailfrom.contains("MicrosoftExchange329e71ec88ae4615bbc36ab6ce41109e@NALCO1.onmicrosoft.com") & mailsubject.contains("Two folders in your mailbox have the same name")
		CASE mailfrom.contains("maverikprod@service-now.com") & mailsubject.contains("dispatch has been CLOSED")
		CASE mailfrom.contains("supportdesk@sprouts.com") & mailsubject.contains("ready to bill")
		CASE mailfrom.contains("support-noreply@avidbill.com") & mailsubject.contains("Your bills have been delivered to Mosaic Management LLC")
		CASE mailfrom.contains("maverikprod@service-now.com") & mailsubject.contains("has been updated")
		CASE mailsubject.contains("has authorized the invoice for Service Request") & mailfrom.contains("provideralert@workordernetwork.corrigo.com")
		CASE mailbody.contains("UPDATE: Status = Closed, Resolution Code = Closed Via Batch.")
		{
			changequeue = true
			targetqueue = "a82590c9-41c9-42b7-9a45-edca9bb87c37"	// CS.US.PEST.EN.CB.EM.WorkOrders
			
			changestatus = true
			targetstatus = "closed"
			debug = 4
		}

		// Brookdale Auto Close
		CASE mailsubject.contains("SCHEDULED Brookdale | Routine PM")
		{
			changequeue = true
			targetqueue = "a82590c9-41c9-42b7-9a45-edca9bb87c37"	// CS.US.PEST.EN.CB.EM.WorkOrders

			sendack = true
			ackmessage = "ACK><br><br>"
			
			changestatus = true
			targetstatus = "closed"
			debug = 5
		}

		// vmUSGrandF.vmPestBilling
		CASE mailfrom.contains("vmUSGrandF.vmPestBilling@ecolab.com")
		{
			changequeue = true
			targetqueue = "1bf62446-1b7e-4404-b0f9-86ec48211d9f"	// CS.US.PEST.EN.CB.EM.Invoicing
			
			priority = 5
			debug = 6
		}

		// Supervisor
		CASE mailfrom.contains("brookdale@ilrs.360facility.net")
		CASE mailfrom.contains("PestElimContractAdminMailbox@ecolab.com")
		CASE mailsubject.contains("Motion sensor activated")
		CASE mailfrom.contains("compass@starbucks.com")
		{
			changequeue = true
			targetqueue = "df5799b1-efd9-4f37-b8d9-f64fe9722d7a"	// CS.US.PEST.EN.CB.EM.Supervisor
			
			priority = 1
			debug = 7
		}

		// Portal
		CASE mailfrom.contains("PESTWOMGMT@ecolab.com") & mailsubject.contains("Rodents and Insects")
		CASE mailfrom.contains("officetrax.noreply@corp.core7.com")
		CASE mailfrom.contains("contact@sprouts.com")
		CASE mailfrom.contains("gei-prod-candidacy@ims-evolve.com")
		CASE mailfrom.contains("panera@service-now.com")
		CASE mailfrom.contains("DOLLAR GENERAL")
		{
			changequeue = true
			targetqueue = "bbfae134-9f40-49d3-8587-e8de40786bfc"	// CS.US.PEST.EN.CB.EM.Portal
			
			priority = 2
			debug = 8
		}

		// From Email - Inv-XferQue
		CASE mailfrom.contains("kelli.peterson@ecolab.com")
		CASE mailfrom.contains("vmUSGrandF.vmPestCredit@ecolab.com")
		CASE mailsubject.contains("Invoice Reprint Request")
		CASE mailsubject.contains("Ecolab Invoice Copies")
		CASE mailsubject.contains("Invoice")
		CASE mailfrom.contains("Kelli.Rother@ecolab.com")
		{
			changequeue = true
			targetqueue = "1bf62446-1b7e-4404-b0f9-86ec48211d9f"	// CS.US.PEST.EN.CB.EM.Invoicing
			
			priority = 5
			debug = 9
		}

		// Workorders
		CASE mailfrom.contains("@servicechannel.com")
		CASE mailsubject.contains("important")
		CASE mailbody.contains("important")
		CASE mailsubject.contains("urgent")
		CASE mailbody.contains("urgent")
		CASE mailsubject.contains("high")
		CASE mailbody.contains("high")
		CASE mailsubject.contains("next day")
		CASE mailbody.contains("next day")
		CASE mailsubject.contains("same day")
		CASE mailbody.contains("same day")
		CASE mailsubject.contains("now")
		CASE mailbody.contains("now")
		CASE mailsubject.contains("immediately")
		CASE mailbody.contains("immediately")
		CASE mailsubject.contains("rush")
		CASE mailbody.contains("rush")
		CASE mailsubject.contains("health department")
		CASE mailbody.contains("health department")
		CASE mailsubject.contains("shut down")
		CASE mailbody.contains("shut down")
		CASE mailsubject.contains("account cancellation/termination")
		CASE mailbody.contains("account cancellation/termination")
		CASE mailsubject.contains("account termination")
		CASE mailbody.contains("account termination")
		CASE mailsubject.contains("cancel my account")
		CASE mailbody.contains("cancel my account")
		CASE mailsubject.contains("infestation")
		CASE mailsubject.contains("customer sighting")
		CASE mailbody.contains("customer sighting")
		CASE mailsubject.contains("no response")
		CASE mailbody.contains("no response")
		CASE mailsubject.contains("multiple messages")
		CASE mailbody.contains("multiple messages")
		CASE mailsubject.contains("multiple calls")
		CASE mailbody.contains("multiple calls")
		CASE mailsubject.contains("caesar")
		CASE mailbody.contains("caesar")
		CASE mailsubject.contains("bally")
		CASE mailbody.contains("bally")
		CASE mailsubject.contains("paris")
		CASE mailbody.contains("paris")
		CASE mailfrom.contains("servicerequest@scalert.com")
		CASE mailfrom.contains("information@verisae.com")
		CASE mailsubject.contains("Pest Product Sale")
		CASE mailfrom.contains("compass@starbucks.com")
		{
			changequeue = true
			targetqueue = "a82590c9-41c9-42b7-9a45-edca9bb87c37"	// CS.US.PEST.EN.CB.EM.WorkOrders
			
			priority = 3
			debug = 10
		}

		// NewLeads
		CASE mailfrom.contains("ecolabecustomerservice@ecolab.com")
		CASE mailfrom.contains("kayla.anderson@ecolab.com")
		{
			changequeue = true
			targetqueue = "9c421b7f-c8bc-4e42-a710-f385d77c8106"	// CS.US.PEST.EN.CB.EM.NewLeads
			
			priority = 3
			debug = 11
		}

		// Weekend Duty
		CASE mailbody.contains("\bwd\b")
		CASE mailsubject.contains("\bwd\b")
		CASE mailbody.contains("weekend duty")
		CASE mailsubject.contains("weekend duty")
		CASE mailbody.contains("weekend duty schedule")
		CASE mailsubject.contains("weekend duty schedule")
		CASE mailbody.contains("vmUSgrandF.vmPestCSCWeekendDuty@ecolab.com")
		{
			changequeue = true
			targetqueue = "e1d0dc96-c057-4ea5-8694-7c73622b823b"	// CS.US.PEST.EN.CB.EM.Wkends
			
			priority = 4
			debug = 12
		}

		// Paging Coverage
		CASE mailsubject.contains("regional meeting")
		CASE mailbody.contains("regional meeting")
		CASE mailsubject.contains("open route")
		CASE mailbody.contains("open route")
		CASE mailsubject.contains("work hours")
		CASE mailbody.contains("work hours")
		CASE mailsubject.contains("paging hours")
		CASE mailbody.contains("paging hours")
		CASE mailsubject.contains("paging assignment")
		CASE mailbody.contains("paging assignment")
		CASE mailsubject.contains("cell phone")
		CASE mailbody.contains("cell phone")
		CASE mailsubject.contains("pagermail")
		CASE mailbody.contains("pagermail")
		{
			changequeue = true
			targetqueue = "a82590c9-41c9-42b7-9a45-edca9bb87c37"	// CS.US.PEST.EN.CB.EM.WorkOrders
			
			priority = 3
			debug = 13
		}

		// Pest Routing
		CASE mailsubject.contains("invoices")
		CASE mailbody.contains("invoices")
		CASE mailsubject.contains("perm transfer")
		CASE mailbody.contains("perm transfer")
		CASE mailsubject.contains("permanent transfer")
		CASE mailbody.contains("permanent transfer")
		CASE mailsubject.contains("temp transfer")
		CASE mailbody.contains("temp transfer")
		CASE mailsubject.contains("temporary transfer")
		CASE mailbody.contains("temporary transfer")
		CASE mailsubject.contains("acct")
		CASE mailbody.contains("acct")
		CASE mailfrom.contains("vmUSGrandF.vmPestBilling@ecolab.com")
		{
			changequeue = true
			targetqueue = "1bf62446-1b7e-4404-b0f9-86ec48211d9f"	// CS.US.PEST.EN.CB.EM.Invoicing
			
			priority = 5
			debug = 14
		}

		// Bulletin Board
		CASE mailsubject.contains("change in effective date")
		CASE mailbody.contains("change in effective date")
		CASE mailsubject.contains("return from work")
		CASE mailbody.contains("return from work")
		CASE mailsubject.contains("empactnx")
		CASE mailbody.contains("empactnx")
		CASE mailsubject.contains("leave of absence")
		CASE mailbody.contains("leave of absence")
		CASE mailsubject.contains("return from leave")
		CASE mailbody.contains("return from leave")
		CASE mailsubject.contains("light duty")
		CASE mailbody.contains("light duty")
		CASE mailsubject.contains("new hire")
		CASE mailbody.contains("new hire")
		CASE mailsubject.contains("number change")
		CASE mailbody.contains("number change")
		CASE mailsubject.contains("updated number transfers")
		CASE mailbody.contains("updated number transfers")
		CASE mailsubject.contains("correction to svsp")
		CASE mailbody.contains("correction to svsp")
		CASE mailsubject.contains("correction to adp")
		CASE mailbody.contains("correction to adp")
		CASE mailfrom.contains("Jaclyn.Brownlee@ecolab.com")
		{
			changequeue = true
			targetqueue = "a82590c9-41c9-42b7-9a45-edca9bb87c37"	// CS.US.PEST.EN.CB.EM.WorkOrders
			
			priority = 3
			debug = 15
		}

		// System Escalation
		CASE mailsubject.contains("US GF")
		CASE mailbody.contains("US GF")
		CASE mailfrom.contains("usrcs@ecolab.com")
		{
			changequeue = true
			targetqueue = "2620d5aa-2e46-42e8-b36e-6cf45530a6fe"	// CS.US.PEST.EN.CB.EM.SysEscltns
			
			priority = 1
			debug = 16
		}
			
		DEFAULT
		{
			changequeue = true
			targetqueue = "a82590c9-41c9-42b7-9a45-edca9bb87c37"	// CS.US.PEST.EN.CB.EM.WorkOrders
			
			priority = 3
			debug = 17
		}
	}
}

// don't send acknowledgements
ASSIGN sendack = false

