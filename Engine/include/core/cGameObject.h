//
// Created by devil on 26.05.17.
//

#ifndef PROJECT_CGAMEOBJECT_H
#define PROJECT_CGAMEOBJECT_H

#include "cFactory.h"
#include "cObject.h"

namespace Core {
    /*   --- Класс игрового объекта, может содержать под объекты ---   */
    class cGameObject : public cFactoryObject, cFactory<cGameObject> {
    protected:
    public:
        cFactory<cObject> *objects;

        cGameObject(std::string name, int id, bool lock = false);
        virtual ~cGameObject();

        virtual void Update(int dt);
    };
}

#endif //PROJECT_CGAMEOBJECT_H
