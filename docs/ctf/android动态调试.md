
低调的天空
1、修改apk配置文件AndroidManifest.xml
     修改配置文件<application>子项，添加 android:debuggable="true" 属性，使apk处于可调式状态，打包回编即可。
           <application android:allowBackup="true" android:debuggable="true" android:icon="@drawable/ic_launcher" android:label="@string/app_name">

2、调试启动apk（两种方式）
     1）使用adb  shell am（没用过）
          使用格式为 【adb shell am start -D -n 包名/包名＋类名】的命令，就可以从控制台启动一个Activity
          该命令的具体用法，参见 adb shell am
     2) 使用手机自带的调试功能（简单方便）
          在手机中的开发者选项中，选择待调试的程序，然后正常单击程序，此时会弹出一个title为“Waiting For Debugger”的对话框，提示用户该应用程序正在等待debugger attach

3、IDA附加app
     在手机中以root身份启动android_server
     在主系统中,使用adb forward tcp:23946 tcp:23946 命令转发端口数据
     完成以上两步后，就可以使用ida附加到应用程序。附加完成后，在ida->Debugger->Debugger Options中，设置event断点，添加 library load/unload 断点，，这样，在app启动之后，每load/unload一个so的时候，都会断住，此时，用户就可以去设置断点，一遍后续调试

4、启动app
     经过以上三步,app处于待调试状态，此时需要启动app，让它跑起来。
     在主系统中，使用ddms查看app的调试端口号（一般为8700），然后使用如下命令：
         jdb -connect com.sun.jdi.SocketAttach:hostname=localhost,port=8700
     此时，app调试状态解除，开始正常运行。
     但是你会发现，“Waiting For Debugger”对话框不见了，但是app还是没有运行，这是因为在ida附加的时候，自动断住了。只要在ida中，按下F9，就可以让app真正的跑起来了。
