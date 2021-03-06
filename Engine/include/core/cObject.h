//
// Created by devil on 26.05.17.
//

#ifndef PROJECT_COBJECT_H
#define PROJECT_COBJECT_H

#include "cFactoryObject.h"

namespace Core {
    /*  --- Базовый класс для составляющих классов игрового объекта ---
     *  --- объектом может быть как звук, так и спрайт, камера или просто трансформация --- */
    class cObject : public cFactoryObject {
    protected:
    public:
        cObject(Paranoia::Engine *engine, std::string name, int id, bool lock = false);

        virtual void Update(int dt);
    };
}

#endif //PROJECT_COBJECT_H
