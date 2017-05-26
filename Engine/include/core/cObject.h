//
// Created by devil on 26.05.17.
//

#ifndef PROJECT_COBJECT_H
#define PROJECT_COBJECT_H

#include "cFactoryObject.h"

namespace Core {
    /*  --- Базовый класс для составляющих классов игрового объекта ---  */
    class cObject : public cFactoryObject {
    protected:
    public:
        cObject(std::string name, int id, bool lock = false);

        void Update(int dt);
    };
}

#endif //PROJECT_COBJECT_H
