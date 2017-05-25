//
// Created by devil on 18.05.17.
//

#include <cstring>
#include "../../include/system/cFile.h"

System::cFile::cFile(std::string name, int id, bool lock) : Core::cFactoryObject(name, id, lock) {
    file = NULL;
    fData = NULL;
}

System::cFile::~cFile() {
    Close();
}

bool System::cFile::Open(FILE_OPEN_TYPE type, bool binaries) {
    if (open)
        Close();

    this->type = type;
    isBinary = binaries;
    std::string modes = "r";

    switch (type) {
        case OPEN_READ: {
            modes = "r";
        }
        break;

        case OPEN_WRITE: {
            modes = "a";
        }
        break;

        case OPEN_WRITE_CLEAR: {
            modes = "w";
        }
        break;

        case OPEN_RW: {
                modes = "a+";
        }
        break;

        case OPEN_RW_CLEAR: {
            modes = "w+";
        }

        default:
            break;
    }

    if (binaries)
        modes += "b";

    file = fopen(name.c_str(), modes.c_str());

    if (file != NULL) {
        open = true;
        return true;
    }

    open = false;
    return false;
}

void System::cFile::Close() {
    if (file != NULL) {
       fclose(file);
        file = NULL;
    }

    ClearData();

    open = false;
}

System::cFileData *System::cFile::Read(long size) {
    if (!open)
        return NULL;

    switch (type) {
        case OPEN_READ:
        case OPEN_RW: {
            ClearData();
            fData = new cFileData();

            if (size == 0) {
                SetPosStart(0);
                fData->size = GetSize();
            } else {
                fData->size = size;
            }

            fData->data = new char[fData->size];

            fData->size = fread(fData->data, isBinary?1:sizeof(char), fData->size, file);
            return fData;
        }
        break;

        case OPEN_WRITE:
        default:
            return NULL;
    }
}

System::cFileData *System::cFile::ReadLine(long size) {
    if (!open)
        return NULL;

    switch (type) {
        case OPEN_READ:
        case OPEN_RW: {
            ClearData();
            fData = new cFileData();

            if (size == 0) {
                SetPosStart(0);
                fData->size = GetSize();
            } else {
                fData->size = size;
            }

            fData->data = new char[fData->size];

            fgets(fData->data, fData->size, file);
            return fData;
        }
            break;

        case OPEN_WRITE:
        default:
            return NULL;
    }
}

System::cFileData *System::cFile::ReadChar() {
    if (!open)
        return NULL;

    switch (type) {
        case OPEN_READ:
        case OPEN_RW: {
            ClearData();
            fData = new cFileData();
            fData->size = 1;
            fData->data = new char;

            *fData->data = (char) getc(file);
            return fData;

        }
            break;

        case OPEN_WRITE:
        default:
            return NULL;
    }
}

void System::cFile::Write(System::cFileData *data) {
    if (!open)
        return;

    switch (type) {
        case OPEN_READ: {
        }
        break;

        case OPEN_WRITE:
        case OPEN_RW: {
            fwrite(data->data, sizeof(char), data->size, file);
        }
        break;

        default:
            break;
    }
}

void System::cFile::Write(std::string message) {
    cFileData fd;

    fd.data = new char[message.size()];
    strcpy(fd.data, message.c_str());
    fd.size = message.size();

    Write(&fd);
}

long System::cFile::GetSize() {
    if (!open)
        return 0;

    long cyr = ftell(file);
    fseek(file, 0, SEEK_END);
    long ret = ftell(file);
    fseek(file, cyr, SEEK_SET);
    return ret;
}

long System::cFile::GetPos() {
    if (!open)
        return 0;

    return ftell(file);
}

void System::cFile::SetPos(long pos) {
    if (!open)
        return;

    fseek(file, pos, SEEK_CUR);
}

void System::cFile::SetPosStart(long pos) {
    if (!open)
        return;

    fseek(file, pos, SEEK_SET);
}

void System::cFile::SetPosEnd(long pos) {
    if (!open)
        return;

    fseek(file, pos, SEEK_END);
}

void System::cFile::ClearData() {
    if (fData) {
        delete fData->data;
        delete fData;
        fData = NULL;
    }
}
