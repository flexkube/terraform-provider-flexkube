package flexkube

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/flexkube/libflexkube/pkg/container/types"
)

func containerConfigMarshal(c types.ContainerConfig) interface{} {
	return []interface{}{
		map[string]interface{}{
			"name":         c.Name,
			"image":        c.Image,
			"privileged":   c.Privileged,
			"args":         c.Args,
			"entrypoint":   c.Entrypoint,
			"port":         portMapMarshal(c.Ports),
			"mount":        mountsMarshal(c.Mounts),
			"network_mode": c.NetworkMode,
			"pid_mode":     c.PidMode,
			"ipc_mode":     c.IpcMode,
			"user":         c.User,
			"group":        c.Group,
		},
	}
}

func containerConfigUnmarshal(i interface{}) types.ContainerConfig {
	j := i.(map[string]interface{})

	args := []string{}
	for _, v := range j["args"].([]interface{}) {
		args = append(args, v.(string))
	}

	entrypoint := []string{}
	for _, v := range j["entrypoint"].([]interface{}) {
		entrypoint = append(entrypoint, v.(string))
	}

	return types.ContainerConfig{
		Name:        j["name"].(string),
		Image:       j["image"].(string),
		Privileged:  j["privileged"].(bool),
		Args:        args,
		Entrypoint:  entrypoint,
		Ports:       portMapUnmarshal(j["port"]),
		Mounts:      mountsUnmarshal(j["mount"]),
		NetworkMode: j["network_mode"].(string),
		PidMode:     j["pid_mode"].(string),
		IpcMode:     j["ipc_mode"].(string),
		User:        j["user"].(string),
		Group:       j["group"].(string),
	}
}

func containerConfigSchema(computed bool) *schema.Schema {
	return requiredBlock(computed, func(computed bool) *schema.Resource {
		return &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name":         requiredString(computed),
				"image":        requiredString(computed),
				"privileged":   optionalBool(computed),
				"args":         optionalStringList(computed),
				"entrypoint":   optionalStringList(computed),
				"port":         portMapSchema(computed),
				"mount":        mountsSchema(computed),
				"network_mode": optionalString(computed),
				"pid_mode":     optionalString(computed),
				"ipc_mode":     optionalString(computed),
				"user":         optionalString(computed),
				"group":        optionalString(computed),
			},
		}
	})
}
