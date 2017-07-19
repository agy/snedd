resource "aws_lambda_function" "snedd-initiator" {
  filename         = "../lambda/initiator/handler.zip"
  function_name    = "snedd-initiator"
  role             = "${aws_iam_role.snedd-initiator.arn}"
  handler          = "handler.Handle"
  source_code_hash = "${base64sha256(file("../lambda/initiator/handler.zip"))}"

  // The golang shim compiles to a python .so file
  runtime = "python2.7"

  environment {
    variables = {
      TTL             = "1800"
      STATEMACHINEARN = "${aws_sfn_state_machine.snedd.id}"
    }
  }
}

resource "aws_lambda_function" "snedd-expirer" {
  filename         = "../lambda/expirer/handler.zip"
  function_name    = "snedd-expirer"
  role             = "${aws_iam_role.snedd-expirer.arn}"
  handler          = "handler.Handle"
  source_code_hash = "${base64sha256(file("../lambda/expirer/handler.zip"))}"

  // The golang shim compiles to a python .so file
  runtime = "python2.7"
}
