//
// Created by devil on 18.05.17.
//

#include "../../include/core/cFactoryObject.h"


Core::cFactoryObject::cFactoryObject(Paranoia::Engine *engine, std::string name, int id, bool lock) {
    this->engine = engine;
    this->name = name;
    this->id = id;
    this->lock = lock;
}

Core::cFactoryObject::~cFactoryObject() {

}

std::string Core::cFactoryObject::GetName() {
    return name;
}

int Core::cFactoryObject::GetId() {
    return id;
}

bool Core::cFactoryObject::GetLock() {
    return lock;
}

void Core::cFactoryObject::Register() {

}
