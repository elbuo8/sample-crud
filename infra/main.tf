module "vpc" {
  source = "github.com/segmentio/stack//vpc"

  name               = "demo"
  environment        = "test"
  cidr               = "10.30.0.0/16"
  external_subnets   = ["10.30.0.0/19"]
  internal_subnets   = ["10.30.64.0/19", "10.30.128.0/19"] # RDS module needs 2
  availability_zones = ["us-west-2a", "us-west-2b"]
}

module "db" {
  source                    = "github.com/segmentio/stack//rds"
  name                      = "samplecrud"                           # No special characters allowed
  password                  = "4c60a640-b4cc-4caa-be5d-fd615bcc12fd"
  vpc_id                    = "${module.vpc.id}"
  subnet_ids                = ["${module.vpc.internal_subnets}"]
  ingress_allow_cidr_blocks = ["${module.vpc.cidr_block}"]
}

resource "aws_security_group" "sample-crud" {
  name        = "sample-app"
  description = "SG "
  vpc_id      = "${module.vpc.id}"

  ingress {
    from_port   = 3000
    to_port     = 3000
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "sample-crud"
    Environment = "test"
  }
}

resource "aws_instance" "sample-crud" {
  ami                    = "ami-093381d21a4fc38d1"                  # ECS AMI
  instance_type          = "t2.micro"
  key_name               = "auth0-laptop"
  subnet_id              = "${module.vpc.external_subnets[0]}"
  vpc_security_group_ids = ["${aws_security_group.sample-crud.id}"]

  # closer SG
  user_data = <<EOF
#!/bin/bash
docker pull elbuo8/sample-crud
docker run -d -p 3000:3000 -e PORT=3000 -e POSTGRES_CONNECTION_DETAILS="${module.db.url}" elbuo8/sample-crud
EOF

  tags {
    Name        = "sample-crud"
    Environment = "test"
  }
}
