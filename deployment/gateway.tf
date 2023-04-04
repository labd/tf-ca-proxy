resource "aws_apigatewayv2_api" "gw" {
  name          = "terraform-registry"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.gw.id
  name        = "$default"
  auto_deploy = true
}


resource "aws_apigatewayv2_route" "service" {
  api_id    = aws_apigatewayv2_api.gw.id
  route_key = "ANY /{proxy+}"
  target    = "integrations/${aws_apigatewayv2_integration.service.id}"
}


resource "aws_apigatewayv2_integration" "service" {
  api_id           = aws_apigatewayv2_api.gw.id
  integration_type = "AWS_PROXY"

  connection_type        = "INTERNET"
  description            = "Integration"
  payload_format_version = "2.0"
  integration_uri        = module.lambda_api.lambda_function_arn
}
