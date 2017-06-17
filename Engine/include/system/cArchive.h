//
// Created by devil on 08.06.17.
//

#ifndef PROJECT_CARCHIVE_H
#define PROJECT_CARCHIVE_H

#include "../core/cFactoryObject.h"
#include "cFile.h"
#include <zlib/zlib.h>

namespace System {
    class cArchive : Core::cFactoryObject {
    protected:
    public:
        cArchive(std::string name, int id, bool lock = false);
    };
}

#endif //PROJECT_CARCHIVE_H
