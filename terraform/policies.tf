data "aws_iam_policy_document" "instance-assume-role-policy" {
  statement {
    actions = [
      "sts:AssumeRole",
    ]

    principals {
      type = "Service"

      identifiers = [
        "ec2.amazonaws.com",
      ]
    }
  }
}

data "aws_iam_policy_document" "ec2-lambda-snedd-initiator" {
  statement {
    actions = [
      "lambda:InvokeFunction",
    ]

    resources = [
      "arn:aws:lambda:${var.region}:${data.aws_caller_identity.current.account_id}:function:snedd-initiator",
    ]
  }
}

resource "aws_iam_policy" "ec2-lambda-snedd-initiator" {
  name   = "ec2-lambda-snedd-initiator"
  policy = "${data.aws_iam_policy_document.ec2-lambda-snedd-initiator.json}"
}

data "aws_iam_policy_document" "sfn-lambda-snedd-expirer" {
  statement {
    actions = [
      "lambda:InvokeFunction",
    ]

    resources = [
      "arn:aws:lambda:${var.region}:${data.aws_caller_identity.current.account_id}:function:snedd-expirer",
    ]
  }
}

resource "aws_iam_policy" "sfn-lambda-snedd-expirer" {
  name   = "sfn-lambda-snedd-expirer"
  policy = "${data.aws_iam_policy_document.sfn-lambda-snedd-expirer.json}"
}

data "aws_iam_policy_document" "lambda-assume-role-policy" {
  statement {
    actions = [
      "sts:AssumeRole",
    ]

    principals {
      type = "Service"

      identifiers = [
        "lambda.amazonaws.com",
      ]
    }
  }
}

data "aws_iam_policy_document" "sfn-assume-role-policy" {
  statement {
    actions = [
      "sts:AssumeRole",
    ]

    principals {
      type = "Service"

      identifiers = [
        "states.${var.region}.amazonaws.com",
      ]
    }
  }
}

data "aws_iam_policy_document" "lambda-exec-snedd-initiator" {
  statement {
    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    resources = [
      "arn:aws:logs:${var.region}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/snedd-initiator:log-stream:*",
    ]
  }

  statement {
    actions = [
      "logs:CreateLogGroup",
    ]

    resources = [
      "arn:aws:logs:${var.region}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/snedd-initiator:*:*",
    ]
  }

  statement {
    actions = [
      "states:StartExecution",
    ]

    resources = [
      "arn:aws:states:${var.region}:${data.aws_caller_identity.current.account_id}:stateMachine:snedd*",
    ]
  }
}

resource "aws_iam_policy" "lambda-exec-snedd-initiator" {
  name   = "lambda-exec-snedd-initiator"
  policy = "${data.aws_iam_policy_document.lambda-exec-snedd-initiator.json}"
}

data "aws_iam_policy_document" "lambda-exec-snedd-expirer" {
  statement {
    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    resources = [
      "arn:aws:logs:${var.region}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/snedd-expirer:log-stream:*",
    ]
  }

  statement {
    actions = [
      "logs:CreateLogGroup",
    ]

    resources = [
      "arn:aws:logs:${var.region}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/snedd-expirer:*:*",
    ]
  }

  statement {
    actions = [
      "ec2:TerminateInstances",
    ]

    resources = [
      "*",
    ]
  }
}

resource "aws_iam_policy" "lambda-exec-snedd-expirer" {
  name   = "lambda-exec-snedd-expirer"
  policy = "${data.aws_iam_policy_document.lambda-exec-snedd-expirer.json}"
}
