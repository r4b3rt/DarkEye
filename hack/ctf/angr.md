#Demo

```
import angr

project = angr.Project("crack", auto_load_libs=False)

@project.hook(0x080486AE)
def print_flag(state):
    print("FLAG SHOULD BE:", state.posix.dumps(0))
    project.terminate_execution()

project.execute()
```
