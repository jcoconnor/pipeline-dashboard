# -*- mode: ruby -*-
# vi: set ft=ruby :

# Useful Vagrant file to setup consistent docker environment for dev building etc.

Vagrant.configure("2") do |config|

  config.vm.box = "genebean/centos-7-docker-ce"
  config.vm.network "forwarded_port", guest: 9292, host: 9292
  config.vm.network "forwarded_port", guest: 16686, host: 16686

  config.vm.provision "shell", inline: <<-SHELL
    curl -L "https://github.com/docker/compose/releases/download/1.25.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
  SHELL

end
