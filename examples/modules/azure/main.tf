
# ======================================
# azure resources
# ======================================
# resource group

data "azurerm_resource_group" "ic" {
  name = var.azure_resource_group_name
}

data "azurerm_express_route_circuit" "ic" {
  name                = var.azure_express_route_circuit_name
  resource_group_name = data.azurerm_resource_group.ic.name
}

# virtual network
resource "azurerm_virtual_network" "vn" {
  name                = "terraform_${var.tag_name}"
  address_space       = ["10.0.0.0/16"]
  location            = data.azurerm_resource_group.ic.location
  resource_group_name = data.azurerm_resource_group.ic.name

  tags = {
    DeployedWith = "terraform"
  }
}

# virtual network subnet
resource "azurerm_subnet" "subnet" {
  name                 = "terraform_${var.tag_name}"
  resource_group_name  = data.azurerm_resource_group.ic.name
  virtual_network_name = azurerm_virtual_network.vn.name
  address_prefix       = "10.0.1.0/24"
}

# subnet gateway
resource "azurerm_subnet" "vngw_subnet" {
  name                 = "GatewaySubnet"
  resource_group_name  = data.azurerm_resource_group.ic.name
  virtual_network_name = azurerm_virtual_network.vn.name
  address_prefix       = "10.0.2.0/24"
}

# public ip for virtual machine
resource "azurerm_public_ip" "ip" {
  name                = "terraform_${var.tag_name}"
  location            = data.azurerm_resource_group.ic.location
  resource_group_name = data.azurerm_resource_group.ic.name
  allocation_method   = "Dynamic"

  tags = {
    DeployedWith = "terraform"
  }
}

# security group allowing ssh from anywhere
resource "azurerm_network_security_group" "nsg" {
  name                = "terraform_${var.tag_name}"
  location            = data.azurerm_resource_group.ic.location
  resource_group_name = data.azurerm_resource_group.ic.name

  security_rule {
    name                       = "Allow_all"
    priority                   = 1001
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "*"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }

  tags = {
    DeployedWith = "terraform"
  }
}

# network interface with public ip
resource "azurerm_network_interface" "nic" {
  name                = "terraform_${var.tag_name}"
  location            = data.azurerm_resource_group.ic.location
  resource_group_name = data.azurerm_resource_group.ic.name

  ip_configuration {
    name                          = "azure-nic"
    subnet_id                     = azurerm_subnet.subnet.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.ip.id
  }

  tags = {
    DeployedWith = "terraform"
  }
}

# virtual machine with network interface
resource "azurerm_virtual_machine" "vm" {
  name                  = "terraform_${var.tag_name}"
  location              = data.azurerm_resource_group.ic.location
  resource_group_name   = data.azurerm_resource_group.ic.name
  network_interface_ids = [azurerm_network_interface.nic.id]
  vm_size               = "Standard_B1ls"

  storage_os_disk {
    name              = "terraform_${var.tag_name}"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }

  os_profile {
    computer_name  = "azure-vm"
    admin_username = "vm-user"
  }

  os_profile_linux_config {
    disable_password_authentication = true
    ssh_keys {
      path     = "/home/vm-user/.ssh/authorized_keys"
      key_data = var.ssh_public_key
    }
  }

  tags = {
    DeployedWith = "terraform"
  }
}

# public ip for virtual gateway
resource "azurerm_public_ip" "vngw_ip" {
  name                = "terraform_${var.tag_name}_vgw"
  location            = data.azurerm_resource_group.ic.location
  resource_group_name = data.azurerm_resource_group.ic.name

  allocation_method = "Dynamic"

  tags = {
    DeployedWith = "terraform"
  }
}

# virtual gateway for express route circuit
resource "azurerm_virtual_network_gateway" "vngw" {
  name                = "terraform_${var.tag_name}"
  location            = data.azurerm_resource_group.ic.location
  resource_group_name = data.azurerm_resource_group.ic.name

  type = "ExpressRoute"
  sku  = "Standard"

  ip_configuration {
    name                          = "terraform_${var.tag_name}"
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.vngw_ip.id
    subnet_id                     = azurerm_subnet.vngw_subnet.id
  }

  tags = {
    DeployedWith = "terraform"
  }
}