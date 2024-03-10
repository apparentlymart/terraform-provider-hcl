terraform {
  required_providers {
    hcl = {
      source = "apparentlymart/hcl"
    }
  }
}

output "result" {
  value = provider::hcl::evalconfig(
    file("${path.module}/config.hcl"),
    "${path.module}/friends-lang.hcldec",
    {},
  )
}
