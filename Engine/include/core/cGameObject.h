//
// Created by devil on 26.05.17.
//

#ifndef PROJECT_CGAMEOBJECT_H
#define PROJECT_CGAMEOBJECT_H

#include "cFactory.h"
#include "cObject.h"

namespace Core {
    /*   --- Класс игрового объекта, может содержать под объекты ---   */
    class cGameObject : public cFactoryObject, cFactory<cObject> {
    protected:
    public:
        cGameObject();
    };
}

#endif //PROJECT_CGAMEOBJECT_H
