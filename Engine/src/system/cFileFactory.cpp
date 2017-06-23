//
// Created by devil on 18.05.17.
//

#include "../../include/system/cFileFactory.h"

bool System::cFileFactory::Open(std::string fName, FILE_OPEN_TYPE type) {
    if (fName == "")
        return false;

    cFile *newFile = FindObject(fName);

    if (newFile == NULL) {
        newFile = new cFile(engine, fName, GetNewID());
    }

    newFile->Open(type);
    return true;

    return false;
}

System::cFileData *System::cFileFactory::Read(std::string fName, unsigned int size) {
    if (fName == "")
        return NULL;

    cFile *newFile = FindObject(fName);

    if (newFile == NULL) {
        return NULL;
    }

    return newFile->Read(size);
}

System::cFileData *System::cFileFactory::ReadLine(std::string fName, unsigned int size) {
    if (fName == "")
        return NULL;

    cFile *newFile = FindObject(fName);

    if (newFile == NULL) {
        return NULL;
    }

    return newFile->ReadLine(size);
}

System::cFileData *System::cFileFactory::ReadChar(std::string fName) {
    if (fName == "")
        return NULL;

    cFile *newFile = FindObject(fName);

    if (newFile == NULL) {
        return NULL;
    }

    return newFile->ReadChar();
}

void System::cFileFactory::SetPos(std::string fName, long pos) {
    if (fName == "")
        return;

    cFile *newFile = FindObject(fName);

    if (newFile == NULL) {
        return;
    }

    newFile->SetPos(pos);
}

void System::cFileFactory::SetPosStart(std::string fName, long pos) {
    if (fName == "")
        return;

    cFile *newFile = FindObject(fName);

    if (newFile == NULL) {
        return;
    }

    newFile->SetPosStart(pos);

}

void System::cFileFactory::SetPosEnd(std::string fName, long pos) {
    if (fName == "")
        return;

    cFile *newFile = FindObject(fName);

    if (newFile == NULL) {
        return;
    }

    newFile->SetPosEnd(pos);

}

void System::cFileFactory::Write(std::string fName, System::cFileData *data) {
    if (fName == "")
        return;

    cFile *newFile = FindObject(fName);

    if (newFile == NULL) {
        return;
    }

    newFile->Write(data);
}

void System::cFileFactory::Write(std::string fName, std::string data) {
    if (fName == "")
        return;

    cFile *newFile = FindObject(fName);

    if (newFile == NULL) {
        return;
    }

    newFile->Write(data);
}

System::cFileFactory::cFileFactory(Paranoia::Engine *engine) : Core::cFactory<cFile>(engine) {

}

System::cFileFactory::~cFileFactory() {
    //Close All
}

