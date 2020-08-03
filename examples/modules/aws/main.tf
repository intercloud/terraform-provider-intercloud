terraform {
  required_version = ">= 0.12"
}

# current user identity
data "aws_caller_identity" "current" {}

# security group allowing ssh from anywhere
resource "aws_security_group" "sg" {
  vpc_id = aws_vpc.vpc.id
  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
  tags = {
    Name         = "terraform#${var.tag_name}"
    DeployedWith = "terraform"
  }
}

# an internet gateway for ou vpc
resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.vpc.id
  tags = {
    Name         = "terraform#${var.tag_name}"
    DeployedWith = "terraform"
  }
}

# dedicated vpc
resource "aws_vpc" "vpc" {
  cidr_block           = "172.30.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = {
    Name         = "terraform#${var.tag_name}"
    DeployedWith = "terraform"
  }
}

# vpc subnet for ec2 instances
resource "aws_subnet" "subnet" {
  vpc_id     = aws_vpc.vpc.id
  cidr_block = "172.30.0.0/24"

  tags = {
    Name         = "terraform#${var.tag_name}"
    DeployedWith = "terraform"
  }
}

# virtual gateway (to be plugged to direct connect)
resource "aws_vpn_gateway" "vpn_gw" {
  vpc_id = aws_vpc.vpc.id
  tags = {
    Name         = "terraform#${var.tag_name}"
    DeployedWith = "terraform"
  }
}

# route table propagating routes from / to virtual gateway
resource "aws_route_table" "rtb" {
  vpc_id = aws_vpc.vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }

  propagating_vgws = [aws_vpn_gateway.vpn_gw.id]

  tags = {
    Name         = "terraform#${var.tag_name}"
    DeployedWith = "terraform"
  }
}

# associate route table to vpc
resource "aws_main_route_table_association" "mrtb" {
  vpc_id         = aws_vpc.vpc.id
  route_table_id = aws_route_table.rtb.id
}

# configure ssh access via public key authentication
resource "aws_key_pair" "key_pair" {
  key_name   = "aws-provider-${var.tag_name}"
  public_key = var.ssh_public_key
}

# ec2 instances with public ip and dns
resource "aws_instance" "vm" {
  depends_on                  = [aws_internet_gateway.igw]
  ami                         = data.aws_ami.latest-ubuntu.id
  key_name                    = aws_key_pair.key_pair.key_name
  associate_public_ip_address = true
  instance_type               = "t2.micro"

  subnet_id       = aws_subnet.subnet.id
  security_groups = [aws_security_group.sg.id]

  tags = {
    Name         = "terraform#${var.tag_name}"
    DeployedWith = "terraform"
  }
}

data "aws_ami" "latest-ubuntu" {
  most_recent = true
  owners      = ["099720109477"] # Canonical

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-*-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}
