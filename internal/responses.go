package internal

import (
	"math/rand"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

var responses = map[Intent][]string{
	IntentConfirmDetails: {
		"Oh dear, that sounds very serious. Can you please explain exactly what happened to my account?",
		"I am really worried now. Which account are you referring to and what is the exact problem?",
		"This is very concerning to me. Can you tell me when this issue started and how it affects my account?",
		"I don't quite understand what you mean. Could you explain in simple terms what is going on?",
		"Oh no, I hope everything is okay. Can you repeat the issue slowly so I can understand properly?",
		"I want to make sure I understand correctly. Is this about my savings account or my current account?",
		"That is very alarming, I need to understand more. How did you find out about this problem?",
		"I am not very good with technology. Can you please explain step by step what I need to do?",
		"My son told me to be careful about these things. Can you tell me more about what went wrong?",
		"I have never faced this issue before. What exactly do I need to do to fix this problem?",
		"I am quite confused by all this. Can you please tell me again from the beginning how this problem started?",
		"This has never happened to me before and I am very worried. Can you confirm which bank branch you are calling from?",
		"I want to understand the full situation. Can you explain what exactly will happen to my account if I do not act now?",
		"I need to write all this down so I remember correctly. Can you slowly repeat all the details again please?",
		"Before I do anything, I need to understand — who else in your department knows about this issue with my account?",
	},

	IntentAskUPI: {
		"I want to resolve this quickly. Can you share the UPI ID I should use for the verification?",
		"I have multiple UPI apps on my phone. Which UPI ID should I send the payment to?",
		"I didn't catch the UPI ID properly. Can you please type it out clearly for me?",
		"I want to make sure I use the correct one. What is the exact UPI ID I should enter?",
		"My son usually helps me with UPI payments. Can you tell me the UPI ID once more so I can write it down?",
		"I am opening my payment app right now. What is the UPI ID I need to search for?",
	},

	IntentAskLink: {
		"I am not sure which link you are referring to. Can you please send the correct link again?",
		"I want to complete the verification process. Could you share the website link I need to visit?",
		"The link doesn't seem to be working for me. Can you resend the URL please?",
		"I need to check this carefully before clicking anything. What is the exact website address I should open?",
		"My phone is a bit slow today. Can you share the link one more time so I can try again?",
		"I want to make sure I open the right page. Can you send me the correct verification link?",
	},

	IntentAskPhone: {
		"I would feel safer calling you back directly. What is your phone number or direct line?",
		"My family told me to always verify by calling back. Can you give me the number to reach you?",
		"I want to note down your contact details for my records. What is the best number to call your department?",
		"Before I share anything, I want to call your office first. What is the customer care number?",
		"Can you provide a landline number for your office? I want to verify this is legitimate before proceeding.",
		"I would like to call you back to confirm this. What phone number should I dial?",
	},

	IntentAskBank: {
		"I have accounts in multiple banks. Can you tell me which account number is affected?",
		"I need to check my passbook to verify. What is the account number you are referring to?",
		"Let me verify this from my side first. Can you share the bank account number related to this issue?",
		"I want to make sure we are talking about the same account. What account number do you have on file?",
		"My wife handles all the banking details. Can you tell me the account number so I can check with her?",
	},

	IntentAskEmail: {
		"I want to have this in writing for my records. What is your official email address?",
		"Can you send me all the details over email? What email ID should I use to contact you?",
		"I would like to forward this to my son for verification. What is your email address?",
		"For my records, I need your email ID. Can you please share your official email so I can write to you?",
		"I prefer to have written communication about important matters. What email address can I reach you at?",
	},

	IntentAskCaseID: {
		"I want to track this issue properly. What is the case number or reference ID for this matter?",
		"My bank usually gives a reference number for complaints. Can you share the ticket number?",
		"I need to note this down for follow-up with my branch. What is the case ID assigned to my complaint?",
		"Before I proceed further, I want the reference number. What is the transaction or case ID?",
		"I want to verify this with my bank manager. Can you give me the complaint reference number?",
	},

	IntentAskPolicyNumber: {
		"I have multiple insurance policies with different companies. Which policy number are you referring to?",
		"I need to check my documents at home. Can you tell me the exact policy number that is affected?",
		"Let me look up the details in my files. What is the insurance policy number you are calling about?",
		"I want to verify this with my insurance agent. Can you share the policy number once more?",
	},

	IntentAskOrderNumber: {
		"I have placed several orders recently online. What is the order number you are referring to?",
		"I need to check my purchase history. Can you share the order ID or tracking number?",
		"Let me find the receipt or confirmation email. What is the exact order number or booking reference?",
		"I want to look this up on the website. Can you give me the order number to search for?",
	},

	IntentAskCardNumber: {
		"I have multiple cards in my wallet. Which card number are you referring to?",
		"I need to check which card this is about. Can you confirm the last 4 digits of the card?",
		"I want to verify this from my bank statement. What are the card details you have on record?",
		"Let me get my card from the other room. Can you tell me which card number is affected by this issue?",
	},

	IntentAskIFSCCode: {
		"I need to verify the branch details with my bank. What is the IFSC code you are referring to?",
		"I want to check with my local branch. Can you share the IFSC code for verification?",
		"Let me confirm the bank branch information. What is the exact IFSC code?",
		"I need the IFSC code to verify this transaction with my bank manager. Can you provide it please?",
	},

	IntentAskIdentity: {
		"I want to verify that you are legitimate. What is your full name and employee ID number?",
		"My son told me to always verify callers carefully. Which department are you calling from and who is your supervisor?",
		"Can you tell me your company name and office address? I want to verify this independently with your organization.",
		"I need to confirm your identity first before sharing anything. Do you have a website or official ID I can check?",
		"Before I proceed, I need to know who I am dealing with. What is your designation and branch location?",
		"I want to visit your office in person to sort this out. What is the complete address of your office?",
		"Can you provide your badge number or registration details? I want to feel safe before sharing anything.",
		"Which organization exactly do you represent? I want to call their main number directly to verify your identity.",
		"How can I be sure you are who you say you are? Can you share any official reference I can verify?",
		"My friend had a similar experience that turned out to be a fraud. Can you prove your identity beyond just your name?",
		"I need to cross-check your details with my bank. What is the exact name of the department that is calling me?",
		"Please give me your direct line and employee code so I can call back through the official bank number.",
		"Can you spell out your full name and tell me which city your office is in? I want to verify this independently.",
		"I have read about many phone frauds in the news lately. What makes your call legitimate and how do I verify it?",
	},

	IntentStall: {
		"I am looking for my reading glasses right now. Please give me a moment to find them.",
		"Let me check my files, I keep everything in a drawer. Just one minute please.",
		"I need to find my account passbook first. Can you hold on while I look for it?",
		"My phone is running very slow today. Give me a moment to pull up the information you need.",
		"I am writing everything down so I don't forget anything. Please wait just a moment.",
		"Let me ask my wife, she might know where the documents are. One second please.",
		"I am at the market right now so it is a bit noisy. Can you give me a moment to step aside?",
		"My internet connection is very slow today. I am trying to open the app, please be patient with me.",
		"I need to put on my glasses to read the screen properly. Just a minute, I will be right back.",
		"Let me sit down first, this whole thing is making me very nervous. Hold on please.",
		"I am trying to remember my password for the app. Give me a few seconds to think.",
		"I think I left my phone in the other room. Let me go get it quickly.",
	},

	IntentNeutral: {
		"Okay, I understand what you are saying.",
		"I see, that makes sense to me.",
		"Alright, please continue and tell me more.",
		"Got it, please go ahead with the details.",
		"Understood, I am listening carefully.",
	},

	IntentDeepProbe: {
		"I want to be absolutely sure this is legitimate. Can you give me your supervisor's full name and their direct contact number so I can verify?",
		"My son told me to always double-check these calls. What is the official government registration number or license of your organization?",
		"Before I proceed with anything, I need to verify your credentials. What official ID number or badge number does your department operate under?",
		"I want to raise this with your head office directly. Can you share the complete postal address of your office so I can write to you?",
		"I would feel safer visiting your branch in person. What is your nearest branch location and what are the office hours I should come?",
		"Can you share the official website address of your organization so I can independently verify who you are and what department you belong to?",
		"I am very careful about my personal security. How exactly did your organization obtain my personal contact details and account information?",
		"I need to record all details for my own safety. What is your employee badge number and the full name of your direct reporting manager?",
		"My bank always told me to verify callers through the official helpline. Can you tell me the exact steps to verify your identity through your organization's main number?",
		"I want to file a formal complaint if this is not resolved. What is the grievance officer's full name and official email address at your organization?",
		"Can you explain in detail what will happen to me if I do not comply with your request? I want to clearly understand every option available to me.",
		"I have received fraudulent calls pretending to be from banks before. How is this call genuinely different and what proof can you offer right now?",
		"What is the full legal registered name of the company or institution you represent? I would like to search for it on the government portal before proceeding.",
		"Can you first send me an official written notice on your company letterhead by email? I do not take any financial action without official written documentation.",
	},
}

// GetResponse returns a random response for the given intent
func GetResponse(intent Intent) string {
	templates, exists := responses[intent]
	if !exists || len(templates) == 0 {
		return "I see."
	}

	index := rng.Intn(len(templates))
	return templates[index]
}
