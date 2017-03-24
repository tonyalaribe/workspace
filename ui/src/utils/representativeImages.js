var fileImageRepresentation = require("../assets/files.png")
var imageImageRepresentation = require("../assets/image.png")
var xlsImageRepresentation = require("../assets/xls.png")
var docImageRepresentation = require("../assets/doc.jpg")



export function GetFileRepresentativeImage(file){

  switch (file.type){
    case "image/png":
      return file.file
    case "image/jpeg":
      return file.file
    case "image/jpg":
      return file.file
    case "image/gif":
      return file.file
    default:
      return fileImageRepresentation
  }
}

export function GetRepresentativeImageByFileExtension(filename){
  let splitFilename = filename.split(".")
  let extension = splitFilename[splitFilename.length-1]
  console.log(extension)
  switch (extension){
    case "png":
      return imageImageRepresentation
    case "jpg":
      return imageImageRepresentation
    case "gif":
      return imageImageRepresentation
    case "jpeg":
      return imageImageRepresentation
    case "xls":
      return xlsImageRepresentation
    case "xlsx":
      return xlsImageRepresentation
    case "doc":
      return docImageRepresentation
    case "docx":
      return docImageRepresentation
    default:
      return fileImageRepresentation
  }
}
