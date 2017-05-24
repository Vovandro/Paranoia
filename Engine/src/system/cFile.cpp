//
// Created by devil on 18.05.17.
//

#include <cstring>
#include "../../include/system/cFile.h"

System::cFile::cFile(std::string name, int id, bool lock) : Core::cFactoryObject(name, id, lock) {
    file = NULL;
}

System::cFile::~cFile() {
    Close();
}

bool System::cFile::Open(FILE_OPEN_TYPE type, bool binaries) {
    if (open)
        Close();

    this->type = type;
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

    open = false;
}

System::cFileData *System::cFile::Read(unsigned int size, bool isLine, bool isWord) {
    if (!open)
        return NULL;

    switch (type) {
        case OPEN_READ: {
            cFileData *ret = new cFileData();
            ret->size = size;
            ret->data = new char[size];

            if (isLine) {
                //f_R->getline(ret->data, ret->size);
                return ret;
            }

            if (isWord) {
               // *f_R >> ret->data;
                return ret;
            }

            //f_R->read(ret->data, ret->size);

            return ret;
        }
        break;

        case OPEN_WRITE: {
            return NULL;
        }
        break;

        case OPEN_RW: {
            cFileData *ret = new cFileData();
            ret->size = size;
            ret->data = new char[size];

            if (isLine) {
                //f_R->getline(ret->data, ret->size);
                return ret;
            }

            if (isWord) {
                //*f_R >> ret->data;
                return ret;
            }

            //f_RW->read(ret->data, ret->size);

            return ret;
        }
        break;

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
