package pkg

import (
	"errors"
	"strings"

	dwv1alpha2 "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	devfileattributes "github.com/devfile/api/v2/pkg/attributes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	kind                                 = "DevWorkspace"
	apiVersion                           = "workspace.devfile.io/v1alpha2"
	nameSuffix                           = "-dw"
	defaultRemoteName                    = "origin"
	defaultDevMemoryLimit                = "8G"
	defaultDevMemoryRequest              = "2G"
	defaultDevCPULimit                   = "4"
	defaultDevCPURequest                 = "1"
	defaultEndpointExposure              = dwv1alpha2.PublicEndpointExposure
	defaultEndpointProtocol              = dwv1alpha2.HTTPEndpointProtocol
	defaultEndpointPath                  = "/"
	defaultEndpointSecure                = false
	defaultDevContainerName              = "cde"
	defaultDevWorkspaceAttributes        = `{"controller.devfile.io/storage-type":"ephemeral","pod-overrides":{"spec":{"shareProcessNamespace":true}}}`
	cheCodeContributionName              = "che-code"
	cheCodeContributionComponentName     = "che-code-runtime-description"
	cheCodeContributionContainerEnvName  = "CODE_HOST"
	cheCodeContributionContainerEnvValue = "0.0.0.0"
	cheCodeContributionURI               = "https://eclipse-che.github.io/che-plugin-registry/main/v3/plugins/che-incubator/che-code/latest/devfile.yaml"
)

var cheCodeContainer = dwv1alpha2.ContainerComponentPluginOverride{
	BaseComponentPluginOverride: dwv1alpha2.BaseComponentPluginOverride{},
	ContainerPluginOverride: dwv1alpha2.ContainerPluginOverride{
		Env: []dwv1alpha2.EnvVarPluginOverride{
			{Name: cheCodeContributionContainerEnvName, Value: cheCodeContributionContainerEnvValue},
		},
	},
}

var cheCodeContribution = dwv1alpha2.ComponentContribution{
	Name: cheCodeContributionName,
	PluginComponent: dwv1alpha2.PluginComponent{
		ImportReference: dwv1alpha2.ImportReference{
			ImportReferenceUnion: dwv1alpha2.ImportReferenceUnion{
				Uri: cheCodeContributionURI,
			},
		},
		PluginOverrides: dwv1alpha2.PluginOverrides{
			Components: []dwv1alpha2.ComponentPluginOverride{
				{
					Name: cheCodeContributionComponentName,
					ComponentUnionPluginOverride: dwv1alpha2.ComponentUnionPluginOverride{
						Container: &cheCodeContainer,
					},
				},
			},
			Commands: nil,
		},
	},
}

func generate(o DebugIDEOptions) (dwv1alpha2.DevWorkspace, error) {
	t, err := template(o)
	if err != nil {
		return dwv1alpha2.DevWorkspace{}, err
	}
	c, err := contribution()
	if err != nil {
		return dwv1alpha2.DevWorkspace{}, err
	}
	d := dwv1alpha2.DevWorkspace{
		TypeMeta: metav1.TypeMeta{
			Kind:       kind,
			APIVersion: apiVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: o.targetPodName + nameSuffix,
		},
		Spec: dwv1alpha2.DevWorkspaceSpec{
			Started:       true,
			Template:      t,
			Contributions: []dwv1alpha2.ComponentContribution{c},
		},
	}
	return d, nil
}

func template(o DebugIDEOptions) (dwv1alpha2.DevWorkspaceTemplateSpec, error) {
	c, err := templateContent(o)
	if err != nil {
		return dwv1alpha2.DevWorkspaceTemplateSpec{}, err
	}
	t := dwv1alpha2.DevWorkspaceTemplateSpec{
		DevWorkspaceTemplateSpecContent: c,
	}
	return t, nil
}

func templateContent(o DebugIDEOptions) (dwv1alpha2.DevWorkspaceTemplateSpecContent, error) {
	// For now only one git repository is supported,
	// but we may support multiple ones in the future
	gitRepos := []string{o.gitRepository}
	dwProjects := make([]dwv1alpha2.Project, 0)
	for _, repo := range gitRepos {
		p, err := project(repo)
		if err != nil {
			return dwv1alpha2.DevWorkspaceTemplateSpecContent{}, err
		}
		dwProjects = append(dwProjects, p)
	}

	// Add the CDE container
	dwComponents := make([]dwv1alpha2.Component, 0)
	c := cdeContainer(o.debugImage)
	dwComponents = append(dwComponents, c)

	// Add the Pod containers
	containers := o.targetPodContainers
	for _, ctr := range containers {
		c := container(ctr)
		dwComponents = append(dwComponents, c)
	}

	// Add the attributes
	dwAttributes, err := attributes()
	if err != nil {
		return dwv1alpha2.DevWorkspaceTemplateSpecContent{}, err
	}

	tc := dwv1alpha2.DevWorkspaceTemplateSpecContent{
		Attributes: dwAttributes,
		Components: dwComponents,
		Projects:   dwProjects,
	}

	return tc, nil
}

func attributes() (devfileattributes.Attributes, error) {
	b := []byte(defaultDevWorkspaceAttributes)
	a := new(devfileattributes.Attributes)
	if err := a.UnmarshalJSON(b); err != nil {
		return devfileattributes.Attributes{}, err
	}
	return *a, nil
}

func project(remote string) (dwv1alpha2.Project, error) {
	p := dwv1alpha2.Project{}
	name, err := projectName(remote)
	if err != nil {
		return p, err
	}
	p.Name = name
	p.ProjectSource = dwv1alpha2.ProjectSource{
		SourceType: "Git",
	}
	g := dwv1alpha2.GitLikeProjectSource{
		CommonProjectSource: dwv1alpha2.CommonProjectSource{},
		Remotes: map[string]string{
			defaultRemoteName: remote,
		},
	}
	gg := dwv1alpha2.GitProjectSource{GitLikeProjectSource: g}
	p.ProjectSource.Git = gg.DeepCopy()
	return p, nil
}

func projectName(remote string) (string, error) {
	remote = strings.TrimSuffix(remote, "/")
	remote = strings.TrimSuffix(remote, ".git")
	i := strings.LastIndex(remote, "/")
	if i == -1 {
		return "", errors.New("Invalid remote name. It doesn't contain the '/': " + remote)
	}
	if i == len(remote)-1 {
		return "", errors.New("invalid remote name. It ends with 2 slashes")
	}
	return remote[i+1:], nil
}

func cdeContainer(image string) dwv1alpha2.Component {
	c := dwv1alpha2.Container{
		Image:         image,
		MemoryLimit:   defaultDevMemoryLimit,
		MemoryRequest: defaultDevMemoryRequest,
		CpuLimit:      defaultDevCPULimit,
		CpuRequest:    defaultDevCPURequest,
	}
	comp := dwv1alpha2.Component{
		Name:       defaultDevContainerName,
		Attributes: nil,
		ComponentUnion: dwv1alpha2.ComponentUnion{
			Container: &dwv1alpha2.ContainerComponent{
				Container: c,
			},
		},
	}
	return comp
}

func container(ctr ContainerInfo) dwv1alpha2.Component {
	vars := make([]dwv1alpha2.EnvVar, 0, len(ctr.env))
	for _, env := range ctr.env {
		v := dwv1alpha2.EnvVar{
			Name:  env.name,
			Value: env.value,
		}
		vars = append(vars, v)
	}
	vols := make([]dwv1alpha2.VolumeMount, 0, len(ctr.volumes))
	for _, vol := range ctr.volumes {
		v := dwv1alpha2.VolumeMount{
			Name: vol.name,
			Path: vol.path,
		}
		vols = append(vols, v)
	}
	ends := make([]dwv1alpha2.Endpoint, 0, len(ctr.endpoints))
	for _, end := range ctr.endpoints {
		secure := defaultEndpointSecure
		e := dwv1alpha2.Endpoint{
			Name:       end.name,
			TargetPort: end.targetPort,
			Exposure:   defaultEndpointExposure,
			Protocol:   defaultEndpointProtocol,
			Secure:     &secure,
			Path:       defaultEndpointPath,
		}
		ends = append(ends, e)
	}

	c := dwv1alpha2.Container{
		Image:         ctr.image,
		Env:           vars,
		VolumeMounts:  vols,
		MemoryLimit:   ctr.memoryLimit,
		MemoryRequest: ctr.memoryRequest,
		CpuLimit:      ctr.cpuLimit,
		CpuRequest:    ctr.cpuRequest,
		Command:       ctr.command,
		Args:          ctr.args,
	}
	comp := dwv1alpha2.Component{
		Name: ctr.name,
		ComponentUnion: dwv1alpha2.ComponentUnion{
			Container: &dwv1alpha2.ContainerComponent{
				Container: c,
				Endpoints: ends,
			},
		},
	}
	return comp
}

func contribution() (dwv1alpha2.ComponentContribution, error) {
	c := cheCodeContribution
	return c, nil
}
