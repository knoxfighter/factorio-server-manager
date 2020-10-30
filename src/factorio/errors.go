package factorio

import "github.com/mroote/factorio-server-manager/bootstrap"

const (
	// credentials
	CredentialsIncomplete = bootstrap.ConstError("Saved Credentials are incomplete")

	// mod_modInfo
	ModInfoInfoJsonNotFoundError = bootstrap.ConstError("ModInfoList: info.json not found in zip-file")

	// factorio api credentials invalid
	FactorioCredentialsInvalid = bootstrap.ConstError("factorio mod portal credentials are invalid")

	// mod_Mods
	ModUploadedModInvalid = bootstrap.ConstError("uploaded mod-file is invalid")

	// saves
	SaveNameInvalid = bootstrap.ConstError("save name invalid")
)

// mod_modInfo
func ModInfoModNotFoundError(data interface{}) error {
	return bootstrap.ErrorWithData{Msg: "ModInfoList: Mod-file not found", Data: data}
}

// mod_modpack
func ModPackAlreadyExistsError(data interface{}) error {
	return bootstrap.ErrorWithData{
		Msg:  "ModPack: Mod pack already exists",
		Data: data,
	}
}
func ModPackDoesNotExistError(data interface{}) error {
	return bootstrap.ErrorWithData{
		Msg:  "ModPack: Mod pack does not exist",
		Data: data,
	}
}

// mod_Mods
func ModPortalError(data interface{}) error {
	return bootstrap.ErrorWithData{
		Msg:  "Mod Portal answered with error",
		Data: data,
	}
}

// mod_modSimple
func ModNotFoundError(data interface{}) error {
	return bootstrap.ErrorWithData{
		Msg:  "Mod not found",
		Data: data,
	}
}

// saves
func SaveNotFoundError(data interface{}) error {
	return bootstrap.ErrorWithData{
		Msg:  "Save not found",
		Data: data,
	}
}
