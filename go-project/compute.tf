resource "aws_instance" "app_server" {
  ami           = "ami-01f23391a59163da9"
  instance_type = "t2.micro"
  key_name      = aws_key_pair.generated_key.key_name

  vpc_security_group_ids = [aws_security_group.app_sg.id]

  tags = {
    Name = "Go-mission"
  }

  user_data = <<-EOF
                #!/bin/bash
                apt-get update
                apt-get upgrade -y

                apt-get install -y apt-transport-https ca-certificates curl software-properties-common
                curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
                echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
                apt-get update
                apt-get install -y docker-ce docker-ce-cli containerd.io
                usermod -aG docker ubuntu
                systemctl enable docker
                systemctl start docker

                sleep 5

                docker run -d -p 8081:8081 --name go-mission dordavidisrael/go-mission:latest

                EOF

  provisioner "local-exec" {
    command = "chmod 400 ${path.module}/private_key.pem"
  }

}


resource "tls_private_key" "key" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "generated_key" {
  key_name   = "my-terraform-key"
  public_key = tls_private_key.key.public_key_openssh
}

resource "local_file" "private_key" {
  content  = tls_private_key.key.private_key_pem
  filename = "${path.module}/private_key.pem"
}

