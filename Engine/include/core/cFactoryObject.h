//
// Created by devil on 18.05.17.
//

#ifndef PROJECT_CFACTORYOBJECT_H
#define PROJECT_CFACTORYOBJECT_H

#include "string"

namespace Core {
    class cFactoryObject {
    protected:
        std::string name;
        int id;
        bool lock;

    public:
        cFactoryObject(std::string name, int id, bool lock = false);
        ~cFactoryObject();

        std::string GetName();
        int GetId();
        bool GetLock();
    };
}

#endif //PROJECT_CFACTORYOBJECT_H
