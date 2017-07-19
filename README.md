# Snedd

The **S**elf (**ne**) **D**estruct **D**evice.

> Snedd has numbers on his forehead because he was condemned to death. [...]
> The penalty for incursion into their Neighbourhood is death by DNA
> expiration. The culprit's DNA is altered so that the body dies exactly one
> year from the date of sentence [...] a few go the whole hog and graft
> display tissue onto the foreheads of executed criminals in the shape of
> digital numbers, to give a read-out of how many days the guy has left.
> Some people think this is unnecessarily bloody-minded, but the Foreheaders
> don't mind too much. Often it gets them served quicker in restuarants
> because the staff can see the guy doesn't have much time to waste.

 -- Michael Marshall Smith, Only Forward

## Overview

While modern infrastructure reduces the need to log in to machines manually
it doesn't (yet) eliminate it. Sometimes you really need to log in to that
one machine to debug that obscure problem. This can eventually result in a
gradual drift in configuration, even when configuration management is used
to mitigate the problem.

I jokingly mentioned to a colleague that it would be cool if a machine was
marked as "tainted" when a user logged in using SSH and that it would
self-destruct after a period of time.

Snedd is the result.

## Demo

[![asciicast](https://asciinema.org/a/DmSOXSVtlKPO2JrbU4PmB6IQe.png)](https://asciinema.org/a/DmSOXSVtlKPO2JrbU4PmB6IQe)

## How it Works

*Note: This has only been tested on Ubuntu 16.04 within AWS*

The node has a custom motd script installed. The script is run on a SSH
login and a motd message is presented to the user. The script runs a Lambda
initiator function which in turn invokes an AWS State Machine which calls an
expirer Lambda function to delete the node after a configurable period of
time.

## Requirements

The following packages are required:
 * A machine with an IAM instance profile allowed to execute the initiator
   Lambda function
 * Two Lambda functions: the initiator and the expirer
 * A Step Function state machine definition

## Build Requirements

All the commands are built with Golang 1.8.3. The Lambda functions use the
[eawsy Lambda shim](https://github.com/eawsy/aws-lambda-go-shim).

## Issues and Caveats

 * Snedd will only work within AWS.
 * Snedd will not trigger on non-interactive SSH logins.
 * If your SSH client uses a control socket (i.e. `ControlPath`) you will
   only be shown the motd on the first login.

## Questions

 * *Q:* Why not just use one Lambda function?
 * *A:* Lambda functions have a maximum execution time of 300 seconds. This
   doesn't give much debugging time.

 * *Q*: How secure is this?
 * *A*: Snedd is a toy system. That said, the motd command retrieves the
   instance's encrypted identity-document and uses this to authorize with
   the initiator Lambda function. The instance ID is retrieved from the
   decrypted document. The validity of the document is checked, however the
   expiry time is not!

 * *Q*: What is the maximum time that can be configured to wait before
   terminating the instance?
 * *A*: This is dependent on AWS Step Functions limits. The current limit is
   one year.

## References and Prior Art

 * [How To Set Up Slack SSH Session Notifications](http://www.ryanbrink.com/slack-ssh-session-notifications/)
 * [Invoke Lambda Function from the CLI](http://docs.aws.amazon.com/lambda/latest/dg/with-userapp-walkthrough-custom-events-invoke.html)
 * [Using AWS Lambda with Scheduled Events](http://docs.aws.amazon.com/lambda/latest/dg/with-scheduled-events.html)

![inspector-gadget-self-destruct](https://cloud.githubusercontent.com/assets/112317/24335641/0ecabbf4-123f-11e7-96f7-8f873c2e1a6c.gif)
