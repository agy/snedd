# Snedd

The Self (ne) Destruct Device.

> Snedd has numbers on his forehead because he was condemned to death. [...]
> The penalty for incursion into their Neighbourhood is death by DNA
> expiration. The culprit's DNA is altered so that the body dies exactly one
> year from the date of sentence [...] a few go the whole hog and graft
> display tissue onto the foreheads of executed criminals in the shape of
> digital numbers, to give a read-out of how many days the guy has left.

 -- Michael Marshall Smith, Only Forward

Unless it's not clear, Snedd isn't meant to be a "serious" solution.

## Overview

While modern infrastructure reduces the need to log in to machines manually
it doesn't (yet) eliminate it. Sometimes you need to log in to that one
machine to debug that obscure problem. This can eventually result in a
gradual drift in configuration, even when Configuration Management is used
to mitigate the problem.

I jokingly mentioned to a colleague that it would be cool if a machine was
marked as "tainted" when a user logged in using SSH and prompted that the
self destruct sequence was initiated.

![inspector-gadget-self-destruct](https://cloud.githubusercontent.com/assets/112317/24335641/0ecabbf4-123f-11e7-96f7-8f873c2e1a6c.gif)

## Problems and Caveats

 * Snedd will not trigger on non-interactive SSH logins.
 * If your SSH client uses a control socket (i.e. `ControlPath`) you will
   only be shown the motd on the first login.

## References and Prior Art

 * [How To Set Up Slack SSH Session Notifications](http://www.ryanbrink.com/slack-ssh-session-notifications/)
 * [Invoke Lambda Function from the CLI](http://docs.aws.amazon.com/lambda/latest/dg/with-userapp-walkthrough-custom-events-invoke.html)
 * [Using AWS Lambda with Scheduled Events](http://docs.aws.amazon.com/lambda/latest/dg/with-scheduled-events.html)
