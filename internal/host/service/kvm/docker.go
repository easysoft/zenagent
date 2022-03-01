package hostAgentKvmService

import "C"
import (
	"context"
	"fmt"
	"github.com/docker/cli/cli/connhelper"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	client "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	agentConf "github.com/easysoft/zv/internal/agent/conf"
	commDomain "github.com/easysoft/zv/internal/comm/domain"
	_commonUtils "github.com/easysoft/zv/internal/pkg/lib/common"
	_logUtils "github.com/easysoft/zv/internal/pkg/lib/log"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

const (
	DockerConnStrRemote = "ssh://%s@%s:22"
)

var (
	DockerCtx    context.Context
	DockerClient *client.Client
)

type DockerService struct {
	existPortMap map[int]bool
}

func NewDockerService() *DockerService {
	s := &DockerService{}
	s.existPortMap = map[int]bool{}
	s.Connect()

	return s
}

func (s *DockerService) ListContainer() (containers []types.Container, err error) {
	containers, err = DockerClient.ContainerList(DockerCtx, types.ContainerListOptions{})
	if err != nil {
		_logUtils.Errorf(err.Error())
	}

	return
}
func (s *DockerService) GetContainer(containerId string) (ret types.Container, err error) {
	containers, err := s.ListContainer()

	for _, container := range containers {
		if container.ID == containerId {
			ret = container
			return
		}
	}

	return
}
func (s *DockerService) GetContainerInfo(containerId string) (ret commDomain.ContainerInfo, err error) {
	contain, err := DockerClient.ContainerInspect(DockerCtx, containerId)
	if err != nil {
		_logUtils.Errorf(err.Error())
		return
	}

	ret.Name = contain.Name
	ret.Image = contain.Image
	if len(contain.HostConfig.PortBindings[nat.Port("20/tcp")]) > 0 {
		sshPort := contain.HostConfig.PortBindings[nat.Port("20/tcp")][0].HostPort
		ret.SshPort, _ = strconv.Atoi(sshPort)
	}

	if len(contain.HostConfig.PortBindings[nat.Port("80/tcp")]) > 0 {
		httpPort := contain.HostConfig.PortBindings[nat.Port("80/tcp")][0].HostPort
		ret.HttpPort, _ = strconv.Atoi(httpPort)
	}

	if len(contain.HostConfig.PortBindings[nat.Port("443/tcp")]) > 0 {
		httpsPort := contain.HostConfig.PortBindings[nat.Port("443/tcp")][0].HostPort
		ret.HttpsPort, _ = strconv.Atoi(httpsPort)
	}

	return
}

func (s *DockerService) CreateContainer(name string, cmd []string) (resp container.ContainerCreateCreatedBody, err error) {
	httpPort := _commonUtils.GetValidPort(58000, 58099, &s.existPortMap)
	sshPort := _commonUtils.GetValidPort(52200, 52299, &s.existPortMap)

	resp, err = DockerClient.ContainerCreate(DockerCtx,
		&container.Config{
			Image: name,
			Cmd:   cmd, //[]string{"echo", "hello world"},
			Env:   []string{"MYSQL_ROOT_PASSWORD=zgxIttknlJK6BpzhWMAz"},
			ExposedPorts: map[nat.Port]struct{}{
				"80/tcp": {},
				"22/tcp": {},
			},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				"22/tcp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: strconv.Itoa(sshPort),
					},
				},
				"80/tcp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: strconv.Itoa(httpPort),
					},
				},
			},
			Mounts: []mount.Mount{
				{
					Source: "/home/aaron/zentaopms",
					Target: "/www/zentaopms",
					Type:   mount.TypeBind,
				},
				{
					Source: "/home/aaron/zentaomysql",
					Target: "/var/lib/mysql",
					Type:   mount.TypeBind,
				},
				{
					Source: "/etc/localtime",
					Target: "/etc/localtime",
					Type:   mount.TypeBind,
				},
			},
		}, nil, nil, "")
	if err != nil {
		_logUtils.Errorf(err.Error())
	}

	return
}
func (s *DockerService) RemoveContainer(containerId string, removeVolumes, removeLinks, force bool) (err error) {
	err = DockerClient.ContainerRemove(DockerCtx, containerId, types.ContainerRemoveOptions{
		RemoveVolumes: removeVolumes,
		RemoveLinks:   removeLinks,
		Force:         force,
	})

	if err != nil {
		_logUtils.Errorf(err.Error())
	}

	return
}

func (s *DockerService) StartContainer(containerId string) (err error) {
	err = DockerClient.ContainerStart(DockerCtx, containerId, types.ContainerStartOptions{})

	if err != nil {
		_logUtils.Errorf(err.Error())
	}

	return
}
func (s *DockerService) StopContainer(containerId string) (err error) {
	err = DockerClient.ContainerStop(DockerCtx, containerId, nil)

	if err != nil {
		_logUtils.Errorf(err.Error())
	}

	return
}
func (s *DockerService) GetContainerLog(containerId string) (ret string, err error) {
	var out io.ReadCloser

	out, err = DockerClient.ContainerLogs(DockerCtx, containerId, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		_logUtils.Errorf(err.Error())
	}

	//stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	data, _ := ioutil.ReadAll(out)
	ret = string(data)

	return
}

func (s *DockerService) ListImage() (images []types.ImageSummary, err error) {
	images, err = DockerClient.ImageList(DockerCtx, types.ImageListOptions{})
	if err != nil {
		_logUtils.Errorf(err.Error())
	}

	return
}

func (s *DockerService) GetImage(imageId string) (ret types.ImageSummary, err error) {
	images, err := s.ListImage()

	for _, image := range images {
		if image.ID == imageId {
			ret = image
			return
		}
	}

	return
}
func (s *DockerService) PullImage(refStr string) (err error) {
	out, err := DockerClient.ImagePull(DockerCtx, refStr, types.ImagePullOptions{})
	if err != nil {
		_logUtils.Errorf(err.Error())
	}

	defer out.Close()
	io.Copy(os.Stdout, out)

	return
}

func (s *DockerService) Connect() {
	var err error
	DockerCtx = context.Background()

	if agentConf.Inst.Host == "" {
		DockerClient, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			_logUtils.Errorf(err.Error())
		}
	} else {
		str := fmt.Sprintf(DockerConnStrRemote, agentConf.Inst.User, agentConf.Inst.Host)
		helper, err := connhelper.GetConnectionHelper(str)

		if err != nil {
			return
		}

		httpClient := &http.Client{
			// No tls
			// No proxy
			Transport: &http.Transport{
				DialContext: helper.Dialer,
			},
		}

		var clientOpts []client.Opt

		clientOpts = append(clientOpts,
			client.WithHTTPClient(httpClient),
			client.WithHost(helper.Host),
			client.WithDialContext(helper.Dialer),
			client.WithAPIVersionNegotiation(),
		)

		DockerClient, err = client.NewClientWithOpts(clientOpts...)

		if err != nil {
			_logUtils.Errorf(err.Error())
		}
	}

	return
}
