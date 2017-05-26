//
// Created by devil on 18.05.17.
//

#ifndef PARANOIA_ENGINE_H
#define PARANOIA_ENGINE_H

#include "system/system.h"
#include "core/core.h"
#include "render/render.h"


enum eStartType
{
    ENGINE_SERVER,
    ENGINE_PC,
};

namespace Paranoia {
    /* Основной класс для работы с движком */
    class Engine {
    protected:
        eStartType type;
        bool run;

    public:
        //Класс дл работы с окном
        System::cWindow *window;
        //Класс для работы с потоками
        System::cThreadFactory *threads;
        //Класс для работы с файлами
        System::cFileFactory *files;
        //Работа с логами
        System::cLog *log;
        //Система рендера
        Render::cRender *render;
        //Система обновления данных
        Core::cUpdate *update;
        //Система состояний
        Core::cStateManager *states;

        Engine(eStartType type);
        ~Engine();

        //Инициализация подсистем движка
        bool Init();
        //Запуск главного цикла движка
        void Start();
        //Отстановка и выгрузка ресурсов
        void Stop();

        //Проверка сообщений
        void handleEvents();
    };
}

#endif //PARANOIA_ENGINE_H
