package internal

import (
	"strings"

	"github.com/Akatana/hqcli/pkg/api"
	"github.com/ktr0731/go-fuzzyfinder"
)

func SelectProject(projects []api.Project) (api.Project, error) {
	idx, err := fuzzyfinder.Find(
		projects,
		func(i int) string {
			return projects[i].Name
		},
	)
	if err != nil {
		return api.Project{}, err
	}
	return projects[idx], nil
}

func SelectTask(tasks []api.Task) (api.Task, error) {
	idx, err := fuzzyfinder.Find(
		tasks,
		func(i int) string {
			return buildTaskPath(tasks[i], tasks)
		},
	)
	if err != nil {
		return api.Task{}, err
	}
	return tasks[idx], nil
}

func buildTaskPath(item api.Task, allItems []api.Task) string {
	idMap := make(map[int]api.Task)
	for _, t := range allItems {
		idMap[t.ID] = t
	}

	var path []string
	current := &item

	for current != nil {
		if current.Name != "" {
			path = append([]string{current.Name}, path...)
		}

		if current.Parent != nil {
			current = current.Parent
		} else if current.ParentID != nil {
			parent, exists := idMap[*current.ParentID]
			if exists {
				current = &parent
			} else {
				break
			}
		} else {
			break
		}
	}

	return strings.Join(path, " / ")
}
