use std::collections::HashMap;
use std::fmt;

#[derive(Debug)]
pub struct DecisionTree {
    pub decision_rules: Vec<DecisionRule>,
}

impl DecisionTree {
    pub fn new() -> Self {
        DecisionTree {
            decision_rules: Vec::new(),
        }
    }

    pub fn train(&mut self, data: &[(&[u8], u8)]) {
        let mut data_map: HashMap<Vec<u8>, u8> = HashMap::new();
        for (input, output) in data {
            let key = input.to_vec();
            let value = *output;
            data_map.insert(key, value);
        }

        self.decision_rules = build_decision_tree(&data_map);
    }

    pub fn predict(&self, input: &[u8]) -> u8 {
        let mut current_node = &self.decision_rules[0];
        while let Some(next_node) = current_node.next_node(input) {
            current_node = &self.decision_rules[next_node as usize];
        }
        current_node.output
    }
}

#[derive(Debug)]
struct DecisionRule {
    input_attribute: String,
    threshold: u8,
    output: u8,
    true_branch: Option<usize>,
    false_branch: Option<usize>,
}

impl DecisionRule {
    fn new(
        input_attribute: String,
        threshold: u8,
        output: u8,
        true_branch: Option<usize>,
        false_branch: Option<usize>,
    ) -> Self {
        DecisionRule {
            input_attribute,
            threshold,
            output,
            true_branch,
            false_branch,
        }
    }

    fn next_node(&self, input: &[u8]) -> Option<usize> {
        if input[self.input_attribute.as_bytes()[0] as usize] as u8 < self.threshold {
            self.true_branch
        } else {
            self.false_branch
        }
    }
}

fn build_decision_tree(data: &HashMap<Vec<u8>, u8>) -> Vec<DecisionRule> {
    // Implement decision tree building algorithm here
}

impl fmt::Display for DecisionTree {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        writeln!(f, "Decision Tree:")?;
        for rule in &self.decision_rules {
            writeln!(
                f,
                "  If {} <= {} then output {}",
                rule.input_attribute, rule.threshold, rule.output
            )?;
        }
        Ok(())
    }
          }
