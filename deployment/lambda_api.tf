
data "aws_region" "current" {}
data "aws_caller_identity" "current" {}



module "lambda_api" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "4.13.0"

  function_name = "terraform-registry"
  description   = "Terraform Registry"
  handler       = "bootstrap"
  runtime       = "provided.al2"
  memory_size   = 128
  timeout       = 5
  publish       = true

  create_package         = false
  local_existing_package = "${path.module}/../dist/terraform-registry.zip"

  create_current_version_allowed_triggers = false

  attach_policies    = false
  number_of_policies = 2
  policies = [
    "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
    "arn:aws:iam::aws:policy/CloudWatchLambdaInsightsExecutionRolePolicy"
  ]

  attach_tracing_policy = true
  tracing_mode          = "Active"

  attach_policy_json = true
  policy_json = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "execute-api:Invoke"
        ],
        Resource = [
          "arn:aws:execute-api:${data.aws_region.current.name}:*:*/*"
        ]
      },
      {
        Effect = "Allow"
        Resource = [
          "arn:aws:codeartifact:eu-west-1:${data.aws_caller_identity.current.account_id}:package/labdigital/terraform-modules/*"
        ]
        Action = [
          "codeartifact:Describe*",
          "codeartifact:Get*",
          "codeartifact:List*",
          "codeartifact:Read*",
        ]
      },

    ]
  })

  environment_variables = {
    "REGISTRY_NAME" : "terraform-modules"
    "REGISTRY_DOMAIN" : "labdigital"
  }

  allowed_triggers = {
    APIGatewayAny = {
      service = "apigateway"
      arn     = aws_apigatewayv2_api.gw.arn
    }
  }

  cloudwatch_logs_retention_in_days = 30
}
