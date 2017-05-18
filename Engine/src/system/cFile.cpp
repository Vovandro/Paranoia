//
// Created by devil on 18.05.17.
//

#include "../../include/system/cFile.h"

System::cFile::cFile(std::string name, int id, bool lock) : Core::cFactoryObject(name, id, lock) {
}

System::cFile::~cFile() {
    Close();
}

bool System::cFile::Open(FILE_OPEN_TYPE type) {
    if (open)
        Close();

    this->type = type;

    switch (type) {
        case OPEN_READ: {


            open = true;
        }
        break;

        case OPEN_WRITE: {


            open = true;
        }
        break;

        case OPEN_RW: {


            open = true;
        }
        break;

        default:
            break;
    }

    return false;
}