// handlers/task_test.go

package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	mock_data "todo-list/mocks/pkg/data"
	data "todo-list/pkg/data"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_data.NewMockTaskRepository(ctrl)
	handler := NewTaskHandler(mockRepo)

	testCases := []struct {
		name       string
		inputTask  data.Task
		expectedID int
	}{
		{
			name:       "Valid Task",
			inputTask:  data.Task{Title: "Sample Task"},
			expectedID: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			taskJSON, err := json.Marshal(tc.inputTask)
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(taskJSON))
			assert.NoError(t, err)

			res := httptest.NewRecorder()

			mockRepo.EXPECT().CreateTask(gomock.Any()).Return(data.Task{ID: tc.expectedID})

			handler.CreateTask(res, req)

			assert.Equal(t, http.StatusCreated, res.Code)

			var createdTask data.Task
			err = json.NewDecoder(res.Body).Decode(&createdTask)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedID, createdTask.ID)
		})
	}
}

func TestUpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_data.NewMockTaskRepository(ctrl)
	handler := NewTaskHandler(mockRepo)

	testCases := []struct {
		name                string
		createTask          data.Task
		updateTaskInput     data.Task
		expectedUpdatedTask data.Task
	}{
		{
			name:                "Valid Task Update",
			createTask:          data.Task{Title: "Task to be Updated"},
			updateTaskInput:     data.Task{ID: 1, Title: "Updated Task", Done: false},
			expectedUpdatedTask: data.Task{ID: 1, Title: "Updated Task", Done: false},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.EXPECT().CreateTask(tc.createTask).Return(data.Task{ID: 1, Title: tc.createTask.Title})

			createTaskJSON, err := json.Marshal(tc.createTask)
			assert.NoError(t, err)

			createReq, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(createTaskJSON))
			assert.NoError(t, err)

			createRes := httptest.NewRecorder()

			handler.CreateTask(createRes, createReq)

			assert.Equal(t, http.StatusCreated, createRes.Code)

			var createdTask data.Task
			err = json.NewDecoder(createRes.Body).Decode(&createdTask)
			assert.NoError(t, err)

			mockRepo.EXPECT().UpdateTask(gomock.Any()).DoAndReturn(func(updatedTask data.Task) error {
				updatedTask.Title = tc.updateTaskInput.Title
				return nil
			})

			updateTaskJSON, err := json.Marshal(tc.updateTaskInput)
			assert.NoError(t, err)

			updateReq, err := http.NewRequest("PUT", "/tasks/"+strconv.Itoa(createdTask.ID), bytes.NewBuffer(updateTaskJSON))
			assert.NoError(t, err)

			updateRes := httptest.NewRecorder()

			handler.UpdateTask(updateRes, updateReq)

			assert.Equal(t, http.StatusOK, updateRes.Code)

			var updatedTask data.Task
			err = json.NewDecoder(updateRes.Body).Decode(&updatedTask)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedUpdatedTask, updatedTask)
		})
	}
}

func TestDeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_data.NewMockTaskRepository(ctrl)
	handler := NewTaskHandler(mockRepo)

	testCases := []struct {
		name        string
		createTask  data.Task
		deleteID    int
		expectedErr error
	}{
		{
			name:        "Valid Task Deletion",
			createTask:  data.Task{ID: 1, Title: "Task to be Deleted"},
			deleteID:    1,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.EXPECT().CreateTask(tc.createTask).Return(data.Task{ID: tc.createTask.ID, Title: tc.createTask.Title})

			createTaskJSON, err := json.Marshal(tc.createTask)
			assert.NoError(t, err)

			createReq, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(createTaskJSON))
			assert.NoError(t, err)

			createRes := httptest.NewRecorder()

			handler.CreateTask(createRes, createReq)

			assert.Equal(t, http.StatusCreated, createRes.Code)

			deleteReq, err := http.NewRequest("DELETE", "/tasks/"+strconv.Itoa(tc.deleteID), nil)
			assert.NoError(t, err)

			deleteRes := httptest.NewRecorder()

			mockRepo.EXPECT().DeleteTask(tc.deleteID).Return(tc.expectedErr)

			handler.DeleteTask(deleteRes, deleteReq)

			if tc.expectedErr != nil {
				assert.Equal(t, http.StatusNotFound, deleteRes.Code)
			} else {
				assert.Equal(t, http.StatusOK, deleteRes.Code)
			}
		})
	}
}
