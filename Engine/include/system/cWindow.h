//
// Created by devil on 18.05.17.
//

#ifndef PROJECT_CWINDOW_H
#define PROJECT_CWINDOW_H

#include <iostream>
#include <SFML/Window.hpp>

namespace Paranoia {
    class Engine;
}

namespace System {
    /* Класс для работы с окном */
    class cWindow {
    protected:
        Paranoia::Engine *engine;
        sf::Window *win;
        int w;
        int h;
        bool isFullscreen;

    public:
        cWindow(Paranoia::Engine *engine);
        ~cWindow();

        /* Инициализация окна
         * isConsole - терминальное окно
         * w,h - ширина и высота окна
         * isFullscreen - полноэкранный режим (не совместим с терминальным окном) */
        bool Init(unsigned int glMajor = 2, unsigned int glMinor = 0,unsigned int antialiasingLevel = 0, bool isConsole = false, int w = 640, int h = 480, bool isFullscreen = false);

        /* Получение дескриптора окна */
        sf::Window *GetWindow();
    };
}

#endif //PROJECT_CWINDOW_H
