resource "aws_sfn_state_machine" "snedd" {
  name     = "snedd"
  role_arn = "${aws_iam_role.snedd-sfn.arn}"

  definition = <<EOF
{
  "Comment": "The Self (ne) Destruct Device",
  "StartAt": "wait",
  "States": {
    "wait": {
      "Type": "Wait",
      "SecondsPath": "$.ttl",
      "Next": "expire_instance"
    },
    "expire_instance": {
      "Type": "Task",
      "Resource": "arn:aws:lambda:${var.region}:${data.aws_caller_identity.current.account_id}:function:snedd-expirer",
      "End": true
    }
  }
}
EOF
}
