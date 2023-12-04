package common

func ValidateArgs(args []string) error {
	if len(args) != 0 {

		return Errorc(InvalidCommandError, "Invalid Usage Of Command Arguments")

	}
	return nil
}

func WriteServiceResponseData(serviceName string, data DiveServiceResponse, cliContext *Cli) error {
	var jsonDataFromFile = Services{}
	err := cliContext.FileHandler().ReadJson("services.json", &jsonDataFromFile)

	if err != nil {
		return err
	}

	_, ok := jsonDataFromFile[serviceName]
	if !ok {
		jsonDataFromFile[serviceName] = &data

	}
	err = cliContext.FileHandler().WriteJson("services.json", jsonDataFromFile)
	if err != nil {
		return err
	}

	return nil
}
