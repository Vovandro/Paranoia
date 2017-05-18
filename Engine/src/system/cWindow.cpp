//
// Created by devil on 18.05.17.
//

#include "../../include/system/cWindow.h"
#include "../../include/engine.h"

System::cWindow::cWindow(class Paranoia::Engine *engine) {
    this->engine = engine;
}


System::cWindow::~cWindow(class Paranoia::Engine *engine) {
}

bool System::cWindow::Init(bool isConsole, int w = 640, int h = 480, bool isFullscreen = false) {
    return true;
}

sf::Window* System::cWindow::GetWindow() {
    return win;
}