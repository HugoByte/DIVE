# NOTE: If you're a VSCode user, you might like our VSCode extension: https://marketplace.visualstudio.com/items?itemName=Kurtosis.kurtosis-extension

NAME_ARG = "name"

# For more information on...
#  - the 'run' function:  https://docs.kurtosis.com/concepts-reference/packages#runnable-packages
#  - the 'plan' object:   https://docs.kurtosis.com/starlark-reference/plan
#  - the 'args' object:   https://docs.kurtosis.com/next/concepts-reference/args
def run(plan, args):
    name = args.get(NAME_ARG, "John Snow")
    plan.print("Hello, " + name)

    # Try out a plan.add_service here (https://docs.kurtosis.com/starlark-reference/plan#add_service)
