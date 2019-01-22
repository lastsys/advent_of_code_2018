#include <cstdio>
#include <cstdlib>
#include <cstdint>
#include <vector>
#include <unordered_set>
#include <iostream>

using namespace std;

typedef vector<int32_t> value_list;

value_list load_file(const char *filename) {
    FILE *stream;
    errno_t err;
    value_list values;

    err = fopen_s(&stream, filename, "r");
    if (err != 0) {
        printf("Unable to open %s.\n", filename);
        exit(EXIT_FAILURE);
    }

    int32_t value;
    while (feof(stream) == 0) {
        fscanf_s(stream, "%d\n", &value);
        values.push_back(value);
    }

    fclose(stream);

    return values;
}

void part1(value_list data) {
    int32_t v = 0;
    for (auto &value : data) {
        v += value;
    }
    printf("Answer part 1: %d\n", v);
}

void part2(value_list data) {
    unordered_set<int32_t> visited;
    int32_t v = 0;
    while (true) {
        for (auto &value : data) {
            v += value;
            if (visited.find(v) != visited.end()) {
                printf("Answer part 2: %d\n", v);
                return;
            } else {
                visited.insert(v);
            }
        }
    }
}

int main() {
    value_list data = load_file("input.txt");
    part1(data);
    part2(data);
    return EXIT_SUCCESS;
}
