//
// Created by devil on 18.05.17.
//

#include "../../include/system/cWindow.h"
#include "../../include/engine.h"

System::cWindow::cWindow(class Paranoia::Engine *engine) {
    this->engine = engine;
}


System::cWindow::~cWindow() {
    win->close();
}

bool System::cWindow::Init(unsigned int glMajor, unsigned int glMinor,unsigned int antialiasingLevel, bool isConsole, int w, int h, bool isFullscreen) {
    if (!isConsole)
    {
        this->w = w;
        this->h = h;
        this->isFullscreen = isFullscreen;

        sf::ContextSettings settings;
        settings.depthBits = 24;
        settings.stencilBits = 8;
        settings.antialiasingLevel = antialiasingLevel;
        settings.majorVersion = glMajor;
        settings.minorVersion = glMinor;

        //Создаем окно
        if (isFullscreen)
            win = new sf::Window(sf::VideoMode(w, h), "Paranoia Engine", sf::Style::Fullscreen, settings);
        else
            win = new sf::Window(sf::VideoMode(w, h), "Paranoia Engine", sf::Style::Default, settings);

        if (win == nullptr)
        {
            std::cout << "Create Window Error: " << std::endl;
            return false;
        }

        win->setVerticalSyncEnabled(true);

        win->setActive(true);
    }

    return true;
}

sf::Window* System::cWindow::GetWindow() {
    return win;
}

void System::cWindow::Update() {
    if (win)
        win->display();
}
