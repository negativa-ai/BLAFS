#include <fstream>
#include <iostream>

/*
Return false if any error happens
*/
bool print_file(std::string file_path)
{
    std::ifstream inf{file_path};
    if (!inf)
    {
        std::cerr << file_path << " not found!" << std::endl;
        return false;
    }

    std::string file_content;
    inf >> file_content;
    std::cout << file_content << std::endl;
    inf.close();
    return true;
}

int main(int argc, char **argv)
{
    if (argc < 2)
    {
        std::cerr << "Please input file path." << std::endl;
        return 1;
    }
    std::string file_path{argv[1]};

    if (!print_file(file_path))
    {
        return 1;
    }

    return 0;
}