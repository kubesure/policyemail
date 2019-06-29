# policyemail

The service generates an pdf and send a email to the insurrer. Policy issued event is picked up by the comms service and put policy pdf and email meta data to S3 and this lambda function processes the file by generating pdf and sending an email. 
