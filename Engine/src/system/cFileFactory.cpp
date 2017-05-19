//
// Created by devil on 18.05.17.
//

#include "../../include/system/cFileFactory.h"

bool System::cFileFactory::Add(std::string fName, FILE_OPEN_TYPE type) {
    if (fName == "")
        return false;

    cFile *newFile = FindObject(fName);

    if (newFile == NULL) {
        newFile = new cFile(fName, GetNewID());
    }

    newFile->Open(type);
    return true;

    return false;
}

System::cFileData *System::cFileFactory::Read(std::string fName, unsigned int size, bool isLine, bool isWord) {
    return NULL;
}
