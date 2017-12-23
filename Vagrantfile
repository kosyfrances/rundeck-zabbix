# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure("2") do |config|

  config.vm.box = "ubuntu/xenial64"
  config.vm.provision "shell", path: "dev/provision.sh"

  # rundeck 2.7.2 server
  config.vm.define "rundeck" do |rundeck|
    rundeck.vm.hostname = "rundeck"
    rundeck.vm.network "private_network", ip: "192.168.44.10"
    rundeck.vm.network "forwarded_port", guest: 4440, host: 4440

    rundeck.vm.provider "virtual box" do |vb|
      vb.memory = 1024
      vb.cpus = 1
    end
  end

  # zabbix 3.4 server
  config.vm.define "zabbix" do |zabbix|
    zabbix.vm.hostname = "zabbix"
    zabbix.vm.network "private_network", ip: "192.168.44.11"

    zabbix.vm.provider "virtual box" do |vb|
      vb.memory = 1024
      vb.cpus = 1
    end
  end

end
