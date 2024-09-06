package ai_models

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/gonum/floats"
	"github.com/gonum/mat"
)

// NeuralNetwork represents a neural network model
type NeuralNetwork struct {
	// Number of inputs
	numInputs int
	// Number of hidden layers
	numHiddenLayers int
	// Number of neurons in each hidden layer
	numNeuronsPerHiddenLayer int
	// Number of outputs
	numOutputs int
	// Weights and biases for each layer
	weights []mat.Matrix
	biases []mat.Vector
}

// NewNeuralNetwork creates a new instance of the neural network
func NewNeuralNetwork(numInputs, numHiddenLayers, numNeuronsPerHiddenLayer, numOutputs int) *NeuralNetwork {
	return &NeuralNetwork{
		numInputs: numInputs,
		numHiddenLayers: numHiddenLayers,
		numNeuronsPerHiddenLayer: numNeuronsPerHiddenLayer,
		numOutputs: numOutputs,
		weights: make([]mat.Matrix, numHiddenLayers+1),
		biases: make([]mat.Vector, numHiddenLayers+1),
	}
}

// Train trains the neural network using the provided training data
func (nn *NeuralNetwork) Train(X, y mat.Matrix) error {
	// Initialize weights and biases randomly
	for i := range nn.weights {
		nn.weights[i] = mat.NewDense(nn.numNeuronsPerHiddenLayer, nn.numInputs, nil)
		nn.biases[i] = mat.NewVecDense(nn.numNeuronsPerHiddenLayer, nil)
		rand.New(rand.NewSource(42)).Normal(nn.weights[i], 0, 1)
		rand.New(rand.NewSource(42)).Normal(nn.biases[i], 0, 1)
	}

	// Train the network using backpropagation
	for iter := 0; iter < 1000; iter++ {
		// Forward pass
		outputs := make([]mat.Matrix, len(X))
		for i, x := range X {
			outputs[i] = nn.forwardPass(x)
		}

		// Calculate error
		error := mat.NewDense(nn.numOutputs, 1, nil)
		for i, y := range y {
			error.Add(y, outputs[i])
			error.Scale(-1)
		}

		// Backward pass
		gradients := make([]mat.Matrix, len(nn.weights))
		for i := range gradients {
			gradients[i] = mat.NewDense(nn.numNeuronsPerHiddenLayer, nn.numInputs, nil)
		}
		for i, x := range X {
			gradients[i].Add(nn.backwardPass(x, error))
		}

		// Update weights and biases
		for i := range nn.weights {
			nn.weights[i].Add(nn.weights[i], gradients[i])
			nn.biases[i].Add(nn.biases[i], gradients[i])
		}
	}

	return nil
}

// forwardPass performs a forward pass through the network
func (nn *NeuralNetwork) forwardPass(x mat.Matrix) mat.Matrix {
	// Calculate output for each layer
	outputs := make([]mat.Matrix, nn.numHiddenLayers+1)
	outputs[0] = x
	for i := 0; i < nn.numHiddenLayers; i++ {
		outputs[i+1] = mat.NewDense(nn.numNeuronsPerHiddenLayer, 1, nil)
		outputs[i+1].Mul(nn.weights[i], outputs[i])
		outputs[i+1].Add(outputs[i+1], nn.biases[i])
		floats.Apply(func(v float64) float64 { return math.Tanh(v) }, outputs[i+1])
	}
	return outputs[nn.numHiddenLayers]
}

// backwardPass performs a backward pass through the network
func (nn *NeuralNetwork) backwardPass(x, error mat.Matrix) mat.Matrix {
	// Calculate error gradients for each layer
	gradients := make([]mat.Matrix, nn.numHiddenLayers+1)
	gradients[nn.numHiddenLayers] = error
	for i := nn.numHiddenLayers - 1; i >= 0; i-- {
		gradients[i] = mat.NewDense(nn.numNeuronsPerHiddenLayer, 1, nil)
		gradients[i].Mul(nn.weights[i].T(), gradients[i+1])
		gradients[i].MulElem(gradients[i], outputs[i])
		gradients[i].Scale(2)
	}
	return gradients[0]
}

// Predict makes a prediction using the trained neural network
func (nn *NeuralNetwork) Predict(x mat.Matrix) mat.Matrix {
	return nn.forwardPass(x)
}
