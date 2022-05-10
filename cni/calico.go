package cni

type Calico struct {
	metadata MetaData
}

func (c Calico) Manifests(template string) string {
	if template == "" {
		template = c.Template()
	}
	if c.metadata.Interface == "" {
		c.metadata.Interface = "interface=" + defaultInterface
	}
	if c.metadata.CIDR == "" {
		c.metadata.CIDR = defaultCIDR
	}

	if c.metadata.CniRepo == "" || c.metadata.CniRepo == defaultCNIRepo {
		c.metadata.CniRepo = "calico"
	}

	if c.metadata.Version == "" {
		c.metadata.Version = "v3.8.2"
	}

	return render(c.metadata, template)
}

func (c Calico) Template() string {
	switch c.metadata.Version {
	case "v3.19.1":
		return CalicoV3191Manifests
	case "v3.8.2":
		return CalicoManifests
	default:
		return CalicoManifests
	}
}
