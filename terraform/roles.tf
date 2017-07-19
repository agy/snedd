# This is the role applied to the EC2 instance
resource "aws_iam_role" "snedd-motd" {
  name               = "snedd-motd"
  assume_role_policy = "${data.aws_iam_policy_document.instance-assume-role-policy.json}"
}

resource "aws_iam_instance_profile" "snedd-motd" {
  name = "snedd-motd"
  role = "${aws_iam_role.snedd-motd.name}"
}

resource "aws_iam_policy_attachment" "snedd-motd" {
  name       = "snedd-motd"
  policy_arn = "${aws_iam_policy.ec2-lambda-snedd-initiator.arn}"

  roles = [
    "${aws_iam_role.snedd-motd.name}",
  ]
}

resource "aws_iam_role" "snedd-initiator" {
  name               = "snedd-initiator"
  assume_role_policy = "${data.aws_iam_policy_document.lambda-assume-role-policy.json}"
}

resource "aws_iam_role_policy_attachment" "snedd-initiator" {
  role       = "${aws_iam_role.snedd-initiator.name}"
  policy_arn = "${aws_iam_policy.lambda-exec-snedd-initiator.arn}"
}

resource "aws_iam_role" "snedd-expirer" {
  name               = "snedd-expirer"
  assume_role_policy = "${data.aws_iam_policy_document.lambda-assume-role-policy.json}"
}

resource "aws_iam_role_policy_attachment" "snedd-expirer" {
  role       = "${aws_iam_role.snedd-expirer.name}"
  policy_arn = "${aws_iam_policy.lambda-exec-snedd-expirer.arn}"
}

resource "aws_iam_role" "snedd-sfn" {
  name               = "snedd-sfn"
  assume_role_policy = "${data.aws_iam_policy_document.sfn-assume-role-policy.json}"
}

resource "aws_iam_role_policy_attachment" "snedd-sfn" {
  role       = "${aws_iam_role.snedd-sfn.name}"
  policy_arn = "${aws_iam_policy.sfn-lambda-snedd-expirer.arn}"
}
