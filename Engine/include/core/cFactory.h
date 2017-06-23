//
// Created by devil on 18.05.17.
//

#ifndef PROJECT_CFACTORY_H
#define PROJECT_CFACTORY_H

#include "cFactoryObject.h"
#include <vector>
#include <typeinfo>

namespace Core {
    /* Базовый класс для реализации фабрик */
    template <class T>
    class cFactory {
    protected:
        //Объекты фабрики
        std::vector<T*> obj;
        //Счетчик ID
        unsigned long ids;
    public:
        cFactory() {ids = 1000;};
        virtual ~cFactory() {};

        // Добавление элемента в фабрику
        virtual void AddObject(T *newObject) {if (newObject == NULL) return; obj.push_back(newObject);};

        virtual T* CreateObject(Paranoia::Engine *engine, std::string name, int id, bool lock = false) {T* tmp = new T(engine, name, id, lock); AddObject(tmp);};

        // Поиск объекта в фабрике
        T* FindObject(std::string name) {for (typename std::vector<T*>::iterator it = obj.begin(); it != obj.end(); ++it) if ((*it)->GetName() == name) return (*it); return NULL;};

        // Удаление объекта
        void RemoveObject(std::string name) {for (typename std::vector<T*>::iterator it = obj.begin(); it != obj.end(); ++it) if ((*it)->GetName() == name){if (!(*it)->GetLock()) obj.erase(it);}};

        // Получение элемента по классу
        template <class R>
        T *Get() {for (int i = 0; i <= obj.size(); i++) if (typeid(*obj[i]).name() == typeid(R).name()) return obj[i]; return NULL;};

        // Получение элемента по имени
        T* Get(std::string name) {for (typename std::vector<T*>::iterator it = obj.begin(); it != obj.end(); ++it) if ((*it)->GetName() == name) return (*it); return NULL;};
        //T* Get(std::string name) {for (int i = 0; i <= obj.size(); i++) if (obj[i]->GetName() == name) return obj[i]; return NULL;};

        unsigned long GetNewID() {return ids++;};

        //Получение всего вектора
        std::vector<T*> *GetAll() {return &obj;};
    };
}

#endif //PROJECT_CFACTORY_H
