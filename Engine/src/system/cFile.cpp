//
// Created by devil on 18.05.17.
//

#include "../../include/system/cFile.h"

System::cFile::cFile(std::string name, int id, bool lock) : Core::cFactoryObject(name, id, lock) {
    f_R = NULL;
    f_W = NULL;
    f_RW = NULL;
}

System::cFile::~cFile() {
    Close();
}

bool System::cFile::Open(FILE_OPEN_TYPE type, bool clear) {
    if (open)
        Close();

    this->type = type;

    switch (type) {
        case OPEN_READ: {
            f_R = new std::ifstream(name.c_str(), std::ios_base::binary);

            if (f_R->is_open()) {
                open = true;
                return true;
            } else {
                open = false;
                return false;
            }
        }
        break;

        case OPEN_WRITE: {
            f_W = new std::ofstream(name.c_str(), clear?(std::ios_base::binary|std::ios_base::trunc):(std::ios_base::binary));

            if (f_W->is_open()) {
                open = true;
                return true;
            } else {
                open = false;
                return false;
            }
        }
        break;

        case OPEN_RW: {
            f_RW = new std::fstream(name.c_str(), clear?(std::ios_base::binary|std::ios_base::trunc):(std::ios_base::binary));

            if (f_RW->is_open()) {
                open = true;
                return true;
            } else {
                open = false;
                return false;
            }
        }
        break;

        default:
            break;
    }

    return false;
}

void System::cFile::Close() {
    if (f_R != NULL) {
        f_R->close();
        delete f_R;
        f_R = NULL;
    }

    if (f_W != NULL) {
        f_W->close();
        delete f_W;
        f_W = NULL;
    }

    if (f_RW != NULL) {
        f_RW->close();
        delete f_RW;
        f_RW = NULL;
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
                f_R->getline(ret->data, ret->size);
                return ret;
            }

            if (isWord) {
                *f_R >> ret->data;
                return ret;
            }

            f_R->read(ret->data, ret->size);

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
                f_R->getline(ret->data, ret->size);
                return ret;
            }

            if (isWord) {
                *f_R >> ret->data;
                return ret;
            }

            f_RW->read(ret->data, ret->size);

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

        case OPEN_WRITE: {
            f_W->write(data->data, data->size);
        }
        break;

        case OPEN_RW: {
            f_RW->write(data->data, data->size);
        }
        break;

        default:
            break;
    }
}
